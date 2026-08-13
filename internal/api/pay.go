package api

import (
	"context"
	"log"
	"net/http"

	"flat-stalker/internal/plan"
	"flat-stalker/internal/repository"

	"github.com/gin-gonic/gin"
)

type InvoiceCreator interface {
	IssuePlanInvoice(ctx context.Context, chatID int64, payload, title, description string, amountKop int) (string, error)
	BotURL() string
}

type PayHandler struct {
	Users    *repository.Users
	Payments *repository.Payments
	Invoices InvoiceCreator
}

type payRequest struct {
	Plan string `json:"plan"`
	Days int    `json:"days"`
}

func (h *PayHandler) Register(r gin.IRoutes) {
	r.POST("/pay", h.Create)
}

func (h *PayHandler) Create(c *gin.Context) {
	chatID, ok := ChatID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req payRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan and days are required"})
		return
	}

	planName := plan.Normalize(req.Plan)
	if planName == plan.Free {
		c.JSON(http.StatusBadRequest, gin.H{"error": "choose plus or pro"})
		return
	}
	kop, ok := plan.PriceKop(planName, req.Days)
	if !ok || kop < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown period"})
		return
	}

	user, err := h.Users.GetByChatID(c.Request.Context(), chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve user"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found, press /start in the bot first"})
		return
	}

	payment, err := h.Payments.Create(c.Request.Context(), user.ID, planName, req.Days, kop)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create payment"})
		return
	}

	title := "FlatStalker " + plan.Label(planName)
	if len(title) > 32 {
		title = plan.Label(planName)
	}
	amount := plan.FormatBYN(kop)
	desc := title + " на " + plan.FormatDaysRU(req.Days) + ". К оплате " + amount + " BYN"
	botURL, err := h.Invoices.IssuePlanInvoice(c.Request.Context(), chatID, payment.Payload, title, desc, kop)
	if err != nil {
		log.Printf("invoice: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to create invoice"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":          true,
		"bot_url":     botURL,
		"invoice_url": "",
		"plan":        planName,
		"days":        req.Days,
		"amount":      amount,
		"amount_kop":  kop,
		"currency":    plan.Currency,
	})
}
