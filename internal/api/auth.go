package api

import (
	"net/http"
	"strings"
	"time"

	"flat-stalker/internal/tgauth"

	"github.com/gin-gonic/gin"
)

const ctxChatIDKey = "tg_chat_id"

// TelegramAuth validates Mini App initData and stores chat id in context.
func TelegramAuth(botToken string, maxAge time.Duration) gin.HandlerFunc {
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
	if raw := strings.TrimSpace(c.GetHeader("X-Telegram-Init-Data")); raw != "" {
		return raw
	}
	auth := strings.TrimSpace(c.GetHeader("Authorization"))
	if strings.HasPrefix(strings.ToLower(auth), "tma ") {
		return strings.TrimSpace(auth[4:])
	}
	return ""
}
