package bot

import (
	"context"
	"log"

	"flat-stalker/internal/config"
	"flat-stalker/internal/repository"

	"github.com/go-telegram/bot"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Bot struct {
	ctx      context.Context
	api      *bot.Bot
	config   *config.Config
	users    *repository.Users
	listings *repository.Listings
}

func New(ctx context.Context, cfg *config.Config, pool *pgxpool.Pool) (*Bot, error) {
	api, err := bot.New(cfg.Telegram.BotToken)
	if err != nil {
		return nil, err
	}

	b := &Bot{
		ctx:      ctx,
		api:      api,
		config:   cfg,
		users:    repository.NewUsers(pool),
		listings: repository.NewListings(pool),
	}
	b.registerHandlers()
	return b, nil
}

func (b *Bot) registerHandlers() {
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeCommand, b.startHandler)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "links", bot.MatchTypeCommand, b.linksHandler)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "help", bot.MatchTypeCommand, b.helpHandler)
}

func (b *Bot) Start() {
	b.api.Start(b.ctx)
}

func (b *Bot) SendMessage(message, chatID string) {
	_, err := b.api.SendMessage(b.ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   message,
	})
	if err != nil {
		log.Println(err)
	}
}
