package api

import (
	"net/http"
	"strconv"

	"flat-stalker/internal/plan"
	"flat-stalker/internal/repository"

	"github.com/gin-gonic/gin"
)

type MeHandler struct {
	Users     *repository.Users
	Intervals plan.Intervals
}

func (h *MeHandler) Register(r gin.IRoutes) {
	r.GET("/api/me", h.Get)
}

func (h *MeHandler) Get(c *gin.Context) {
	chatID, err := strconv.ParseInt(c.Query("chat_id"), 10, 64)
	if err != nil || chatID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chat_id is required"})
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

	userPlan := plan.Normalize(user.Plan)
	interval := h.Intervals.For(userPlan)

	c.JSON(http.StatusOK, gin.H{
		"ok":          true,
		"plan":        userPlan,
		"plan_label":  plan.Label(userPlan),
		"interval":    interval.String(),
		"interval_ms": interval.Milliseconds(),
		"intervals": gin.H{
			plan.Free: h.Intervals.Free.String(),
			plan.Plus: h.Intervals.Plus.String(),
			plan.Pro:  h.Intervals.Pro.String(),
		},
	})
}
