package bot

import (
	"context"
	"fmt"
	"log"
	"strings"

	"flat-stalker/internal/plan"

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
	listings, err := b.listings.ListByChatID(ctx, chatID)
	if err != nil {
		log.Printf("links: chat_id=%d: %v", chatID, err)
		_, _ = b.api.SendMessage(b.ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Не удалось получить ссылки. Попробуй позже.",
		})
		return
	}

	if len(listings) == 0 {
		_, _ = b.api.SendMessage(b.ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Пока нет сохранённых ссылок. Добавь ссылку поиска Kufar в кабинете Mini App.",
		})
		return
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Твои ссылки (%d):\n\n", len(listings))
	for i, l := range listings {
		status := "активна"
		if l.Paused {
			status = "пауза"
		}
		fmt.Fprintf(&sb, "%d. [%s] %s\n", i+1, status, l.URL)
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

func (b *Bot) statusHandler(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	user, err := b.users.GetByChatID(ctx, chatID)
	if err != nil {
		log.Printf("status: chat_id=%d: %v", chatID, err)
		_, _ = b.api.SendMessage(b.ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Не удалось получить статус. Попробуй позже.",
		})
		return
	}
	if user == nil {
		_, _ = b.api.SendMessage(b.ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Сначала нажми /start.",
		})
		return
	}

	userPlan := plan.Normalize(user.Plan)
	interval := plan.Intervals{
		Free: b.config.Worker.Free,
		Plus: b.config.Worker.Plus,
		Pro:  b.config.Worker.Pro,
	}.For(userPlan)

	text := fmt.Sprintf(
		"Твой тариф: %s\nПроверка объявлений: %s",
		plan.Label(userPlan),
		plan.FormatIntervalRU(interval),
	)

	_, err = b.api.SendMessage(b.ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
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
		Text:      fmt.Sprintf("Команды: /start, /links, /status, /help\nВопросы: <b>%s</b>", b.config.Telegram.SupportContact),
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		log.Println(err)
	}
}

func botBool(v bool) *bool {
	return new(v)
}
