package api

import (
	"net/http"
	"strings"
	"time"

	"flat-stalker/internal/repository"
	"flat-stalker/internal/tgauth"

	"github.com/gin-gonic/gin"
)

const ctxChatIDKey = "tg_chat_id"

// TelegramAuth validates Mini App initData and stores chat id in context.
func TelegramAuth(botToken string, maxAge time.Duration, bans *repository.Bans, support string) gin.HandlerFunc {
	return func(c *gin.Context) {
		initData := extractInitData(c)
		if initData == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "telegram auth required"})
			return
		}
		if err := tgauth.Validate(initData, botToken, maxAge); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid telegram auth"})
			return
		}
		user, err := tgauth.ParseUser(initData)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid telegram user"})
			return
		}
		ban, err := bans.Get(c.Request.Context(), user.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "auth failed"})
			return
		}
		if ban != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "access closed",
				"code":    "banned",
				"support": support,
				"reason":  ban.Reason,
			})
			return
		}
		c.Set(ctxChatIDKey, user.ID)
		c.Next()
	}
}

func ChatID(c *gin.Context) (int64, bool) {
	v, ok := c.Get(ctxChatIDKey)
	if !ok {
		return 0, false
	}
	id, ok := v.(int64)
	return id, ok && id != 0
}

func extractInitData(c *gin.Context) string {
	if h := strings.TrimSpace(c.GetHeader("X-Telegram-Init-Data")); h != "" {
		return h
	}
	auth := strings.TrimSpace(c.GetHeader("Authorization"))
	if strings.HasPrefix(strings.ToLower(auth), "tma ") {
		return strings.TrimSpace(auth[4:])
	}
	return ""
}
