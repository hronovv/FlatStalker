package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	godotenv.Load(".env")
	token := os.Getenv("BOT_TOKEN")
	b, err := bot.New(token)
	if err != nil {
		panic(err)
	}

	TgBot := Bot{
		b:   b,
		ctx: ctx,
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/message", TgBot.Message)

	go func() {
		b.Start(ctx)
	}()

	router.Run(":8080")
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}

type Bot struct {
	ctx context.Context
	b   *bot.Bot
}

type MessageDTO struct {
	Msg    string `json:"message"`
	ChatID string `json:"chat_id"`
}

func (TgBot *Bot) Message(c *gin.Context) {
	var message MessageDTO
	if err := c.BindJSON(&message); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	TgBot.SendMessage(message.Msg, message.ChatID)
}

func (TgBot *Bot) SendMessage(message, chat_id string) {
	TgBot.b.SendMessage(TgBot.ctx, &bot.SendMessageParams{
		ChatID: chat_id,
		Text:   message,
	})
}
