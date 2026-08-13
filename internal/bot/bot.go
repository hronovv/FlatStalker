package bot

import (
	"context"

	"flat-stalker/internal/config"
	"flat-stalker/internal/repository"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Bot struct {
	ctx      context.Context
	api      *bot.Bot
	config   *config.Config
	users    *repository.Users
	listings *repository.Listings
	bans     *repository.Bans
	payments *repository.Payments
	username string
}

func New(ctx context.Context, cfg *config.Config, users *repository.Users, listings *repository.Listings, bans *repository.Bans, payments *repository.Payments) (*Bot, error) {
	api, err := bot.New(cfg.Telegram.BotToken, bot.WithDefaultHandler(func(context.Context, *bot.Bot, *models.Update) {}))
	if err != nil {
		return nil, err
	}

	me, err := api.GetMe(ctx)
	if err != nil {
		return nil, err
	}

	b := &Bot{
		ctx:      ctx,
		api:      api,
		config:   cfg,
		users:    users,
		listings: listings,
		bans:     bans,
		payments: payments,
		username: me.Username,
	}
	b.registerHandlers()
	return b, nil
}

func (b *Bot) registerHandlers() {
	b.api.RegisterHandlerMatchFunc(func(u *models.Update) bool {
		return u.PreCheckoutQuery != nil
	}, b.preCheckoutHandler)
	b.api.RegisterHandlerMatchFunc(func(u *models.Update) bool {
		return u.Message != nil && u.Message.SuccessfulPayment != nil
	}, b.successfulPaymentHandler)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeCommand, b.startHandler)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "links", bot.MatchTypeCommand, b.linksHandler)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "status", bot.MatchTypeCommand, b.statusHandler)
	b.api.RegisterHandler(bot.HandlerTypeMessageText, "help", bot.MatchTypeCommand, b.helpHandler)
	b.api.RegisterHandler(bot.HandlerTypeCallbackQueryData, linksCallbackPrefix, bot.MatchTypePrefix, b.linksCallbackHandler)
}

func (b *Bot) Start() {
	b.api.Start(b.ctx)
}

func (b *Bot) NotifyAd(ctx context.Context, chatID int64, text string) error {
	_, err := b.api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
	return err
}
