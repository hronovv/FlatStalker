package api

import (
	"net/http"

	"flat-stalker/internal/plan"
	"flat-stalker/internal/repository"

	"github.com/gin-gonic/gin"
)

type MeHandler struct {
	Users     *repository.Users
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
