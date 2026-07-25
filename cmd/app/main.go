package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"flat-stalker/internal/config"
	"flat-stalker/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func main() {
	cfg := config.MustLoad()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	pool := db.MustNewPool(ctx, cfg.Database)
	defer pool.Close()
	log.Println("database connected")

	b, err := bot.New(cfg.Telegram.BotToken)
	if err != nil {
		panic(err)
	}

	tgBot := Bot{
		b:      b,
		ctx:    ctx,
		config: cfg,
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "help", bot.MatchTypeCommand, tgBot.helpHandler)

	gin.SetMode(cfg.App.GinMode)
	router := gin.Default()
	router.POST("/message", tgBot.Message)

	go func() {
		b.Start(ctx)
	}()

	log.Printf("http listening on %s", cfg.Server.Addr)
	if err := router.Run(cfg.Server.Addr); err != nil {
		log.Fatal(err)
	}
}

type Bot struct {
	ctx    context.Context
	b      *bot.Bot
	config *config.Config
}

type MessageDTO struct {
	Msg    string `json:"message"`
	ChatID string `json:"chat_id"`
}

func (tgBot *Bot) Message(c *gin.Context) {
	var message MessageDTO
	if err := c.BindJSON(&message); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	tgBot.SendMessage(message.Msg, message.ChatID)
}

func (tgBot *Bot) helpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := tgBot.b.SendMessage(tgBot.ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      fmt.Sprintf("По всем вопросам пишите сюда <b>%s</b>", tgBot.config.Telegram.SupportContact),
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		log.Println(err)
	}
}

func (tgBot *Bot) SendMessage(message, chatID string) {
	_, err := tgBot.b.SendMessage(tgBot.ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   message,
	})
	if err != nil {
		log.Println(err)
	}
}
