package api

import (
	"errors"
	"net/http"
	"time"

	"flat-stalker/internal/models"
	"flat-stalker/internal/plan"
	"flat-stalker/internal/repository"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type MeHandler struct {
	Users     *repository.Users
	Listings  *repository.Listings
	Intervals plan.Intervals
}

func (h *MeHandler) Register(r gin.IRoutes) {
	r.GET("/me", h.Get)
}

func (h *MeHandler) Get(c *gin.Context) {
	chatID, ok := ChatID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var (
		user     *models.User
		listings []models.Listing
	)

	g, ctx := errgroup.WithContext(c.Request.Context())
	g.Go(func() error {
		var err error
		user, err = h.Users.GetByChatID(ctx, chatID)
		return err
	})
	g.Go(func() error {
		var err error
		listings, err = h.Listings.ListByChatID(ctx, chatID)
		return err
	})
	if err := g.Wait(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load cabinet"})
		return
	}
	if user == nil {
		created, err := h.Users.CreateByChatID(c.Request.Context(), chatID)
		if errors.Is(err, repository.ErrBanned) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access closed", "code": "banned"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load cabinet"})
			return
		}
		user = created
	}

	userPlan := plan.Effective(user.Plan, user.PlanExpiresAt, time.Now())
	interval := h.Intervals.For(userPlan)
	limit := plan.LinkLimit(userPlan)

	resp := gin.H{
		"ok":          true,
		"plan":        userPlan,
		"plan_label":  plan.Label(userPlan),
		"interval":    interval.String(),
		"interval_ms": interval.Milliseconds(),
		"link_limit":  limit,
		"intervals": gin.H{
			plan.Free: h.Intervals.Free.String(),
			plan.Plus: h.Intervals.Plus.String(),
			plan.Pro:  h.Intervals.Pro.String(),
		},
		"link_limits": gin.H{
			plan.Free: plan.LinkLimit(plan.Free),
			plan.Plus: plan.LinkLimit(plan.Plus),
			plan.Pro:  plan.LinkLimit(plan.Pro),
		},
		"prices": plan.PriceCatalog(),
		"links":  listingsJSON(listings),
	}
	if user.PlanExpiresAt != nil {
		resp["plan_expires_at"] = user.PlanExpiresAt.UTC().Format("2006-01-02T15:04:05Z")
	}
	c.JSON(http.StatusOK, resp)
}
