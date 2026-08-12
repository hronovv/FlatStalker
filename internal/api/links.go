package api

import (
	"net/http"
	"strconv"
	"strings"

	"flat-stalker/internal/models"
	"flat-stalker/internal/plan"
	"flat-stalker/internal/repository"
	"flat-stalker/internal/source/kufar"
	"flat-stalker/internal/worker"

	"github.com/gin-gonic/gin"
)

type LinksHandler struct {
	Users    *repository.Users
	Listings *repository.Listings
	Seen     *repository.SeenAds
	Kufar    *kufar.Client
}

type addLinkRequest struct {
	URL string `json:"url" binding:"required"`
}

type pauseLinkRequest struct {
	Paused *bool `json:"paused" binding:"required"`
}

func (h *LinksHandler) Register(r gin.IRoutes) {
	r.POST("/links", h.Add)
	r.PATCH("/links/:id", h.SetPaused)
	r.DELETE("/links/:id", h.Delete)
}

func listingsJSON(listings []models.Listing) []gin.H {
	items := make([]gin.H, 0, len(listings))
	for _, l := range listings {
		items = append(items, listingJSON(l))
	}
	return items
}

func listingJSON(l models.Listing) gin.H {
	return gin.H{
		"id":     l.ID,
		"url":    l.URL,
		"paused": l.Paused,
	}
}

func (h *LinksHandler) Add(c *gin.Context) {
	chatID, ok := ChatID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req addLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
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

	user, err := h.Users.GetByChatID(c.Request.Context(), chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve user"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found, press /start in the bot first"})
		return
	}

	limit := plan.LinkLimit(user.Plan)
	count, err := h.Listings.CountByUserID(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save link"})
		return
	}
	if count >= limit {
		c.JSON(http.StatusConflict, gin.H{
			"error": "link limit reached",
			"code":  "link_limit",
			"limit": limit,
			"used":  count,
		})
		return
	}

	listing, created, err := h.Listings.Add(c.Request.Context(), user.ID, chatID, link)
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

	worker.SeedFromURLAsync(h.Kufar, h.Seen, listing.ID, listing.URL)

	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"created": true,
		"id":      listing.ID,
		"url":     listing.URL,
		"paused":  listing.Paused,
	})
}

func (h *LinksHandler) SetPaused(c *gin.Context) {
	chatID, ok := ChatID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	listingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || listingID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid link id"})
		return
	}

	var req pauseLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Paused == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "paused is required"})
		return
	}

	listing, err := h.Listings.SetPaused(c.Request.Context(), listingID, chatID, *req.Paused)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update link"})
		return
	}
	if listing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "link not found"})
		return
	}

	c.JSON(http.StatusOK, listingJSON(*listing))
}

func (h *LinksHandler) Delete(c *gin.Context) {
	chatID, ok := ChatID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	listingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || listingID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid link id"})
		return
	}

	okDel, err := h.Listings.Delete(c.Request.Context(), listingID, chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete link"})
		return
	}
	if !okDel {
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

		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
