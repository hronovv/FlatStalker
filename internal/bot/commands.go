package bot

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *Bot) startHandler(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	if _, err := b.users.CreateByChatID(ctx, chatID); err != nil {
		log.Printf("start: create user chat_id=%d: %v", chatID, err)
	}

	_, err := b.api.SendMessage(b.ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Привет! Ты в FlatStalker. Открой кабинет и добавь ссылку поиска аренды с Kufar. Команда /links покажет сохранённые.",
	})
	if err != nil {
		log.Println(err)
	}
}

func (b *Bot) linksHandler(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	urls, err := b.listings.ListURLsByChatID(ctx, chatID)
	if err != nil {
		log.Printf("links: chat_id=%d: %v", chatID, err)
		_, _ = b.api.SendMessage(b.ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Не удалось получить ссылки. Попробуй позже.",
		})
		return
	}

	if len(urls) == 0 {
		_, _ = b.api.SendMessage(b.ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Пока нет сохранённых ссылок. Добавь ссылку поиска Kufar в кабинете Mini App.",
		})
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Твои ссылки (%d):\n\n", len(urls)))
	for i, u := range urls {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, u))
	}

	_, err = b.api.SendMessage(b.ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   sb.String(),
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: botBool(true),
		},
	})
	if err != nil {
		log.Println(err)
	}
}

func (b *Bot) helpHandler(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	_, err := b.api.SendMessage(b.ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      fmt.Sprintf("Команды: /start, /links, /help\nВопросы: <b>%s</b>", b.config.Telegram.SupportContact),
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		log.Println(err)
	}
}

func botBool(v bool) *bool {
	return new(v)
}
