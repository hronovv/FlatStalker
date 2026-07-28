package api

import (
	"net/http"
	"strconv"
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

type pauseLinkRequest struct {
	ChatID int64 `json:"chat_id" binding:"required"`
	Paused *bool `json:"paused" binding:"required"`
}

type deleteLinkRequest struct {
	ChatID int64 `json:"chat_id" binding:"required"`
}

func (h *LinksHandler) Register(r gin.IRoutes) {
	r.GET("/api/links", h.List)
	r.POST("/api/links", h.Add)
	r.PATCH("/api/links/:id", h.SetPaused)
	r.DELETE("/api/links/:id", h.Delete)
}

func (h *LinksHandler) List(c *gin.Context) {
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

	listings, err := h.Listings.ListByChatID(c.Request.Context(), chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list links"})
		return
	}

	items := make([]gin.H, 0, len(listings))
	for _, l := range listings {
		items = append(items, gin.H{
			"id":     l.ID,
			"url":    l.URL,
			"paused": l.Paused,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":    true,
		"links": items,
	})
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
		"paused":  listing.Paused,
	})
}

func (h *LinksHandler) SetPaused(c *gin.Context) {
	listingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || listingID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid link id"})
		return
	}

	var req pauseLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Paused == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chat_id and paused are required"})
		return
	}

	listing, err := h.Listings.SetPaused(c.Request.Context(), listingID, req.ChatID, *req.Paused)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update link"})
		return
	}
	if listing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "link not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":     true,
		"id":     listing.ID,
		"url":    listing.URL,
		"paused": listing.Paused,
	})
}

func (h *LinksHandler) Delete(c *gin.Context) {
	listingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || listingID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid link id"})
		return
	}

	var req deleteLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.ChatID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chat_id is required"})
		return
	}

	ok, err := h.Listings.Delete(c.Request.Context(), listingID, req.ChatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete link"})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "link not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "deleted": true})
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
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
