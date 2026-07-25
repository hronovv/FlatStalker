package bot

import (
	"github.com/gin-gonic/gin"
)

type MessageDTO struct {
	Msg    string `json:"message"`
	ChatID string `json:"chat_id"`
}

func (b *Bot) Message(c *gin.Context) {
	var message MessageDTO
	if err := c.BindJSON(&message); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	b.SendMessage(message.Msg, message.ChatID)
}
