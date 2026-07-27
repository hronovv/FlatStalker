package api

import (
	"net/http"
	"strings"

	"flat-stalker/internal/repository"
	"flat-stalker/internal/source/kufar"

	"github.com/gin-gonic/gin"
)

type LinksHandler struct {
	Users    *repository.Users
	Listings *repository.Listings
}

type addLinkRequest struct {
	ChatID int64  `json:"chat_id" binding:"required"`
	URL    string `json:"url" binding:"required"`
}

func (h *LinksHandler) Register(r gin.IRoutes) {
	r.POST("/api/links", h.Add)
}

func (h *LinksHandler) Add(c *gin.Context) {
	var req addLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chat_id and url are required"})
		return
	}

	link := strings.TrimSpace(req.URL)
	if link == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}
	if err := kufar.ValidateSearchURL(link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Users.GetByChatID(c.Request.Context(), req.ChatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve user"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found, press /start in the bot first"})
		return
	}

	listing, created, err := h.Listings.Add(c.Request.Context(), user.ID, link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save link"})
		return
	}
	if !created {
		c.JSON(http.StatusOK, gin.H{
			"ok":      true,
			"created": false,
			"message": "link already exists",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"created": true,
		"id":      listing.ID,
		"url":     listing.URL,
	})
}

func CORS(allowedOrigins []string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	allowAll := false
	for _, origin := range allowedOrigins {
		origin = strings.TrimSpace(origin)
		if origin == "*" {
			allowAll = true
			continue
		}
		if origin != "" {
			allowed[origin] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if allowAll && origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
		} else if _, ok := allowed[origin]; ok {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
		}

		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
