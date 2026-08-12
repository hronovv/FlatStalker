package bot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	appmodels "flat-stalker/internal/models"
	"flat-stalker/internal/plan"
	"flat-stalker/internal/repository"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const linksCallbackPrefix = "lnk:"

func (b *Bot) startHandler(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	if b.rejectIfBanned(ctx, chatID) {
		return
	}
	if _, err := b.users.CreateByChatID(ctx, chatID); err != nil {
		if errors.Is(err, repository.ErrBanned) {
			reason := ""
			if ban, banErr := b.bans.Get(ctx, chatID); banErr == nil && ban != nil {
				reason = ban.Reason
			}
			b.sendBanNotice(ctx, chatID, reason)
			return
		}
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
	if b.rejectIfBanned(ctx, chatID) {
		return
	}
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

	_, err = b.api.SendMessage(b.ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   fmt.Sprintf("Твои ссылки (%d). Управление кнопками ниже:", len(listings)),
	})
	if err != nil {
		log.Println(err)
		return
	}

	for i, listing := range listings {
		_, err = b.api.SendMessage(b.ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   formatListingMessage(i+1, listing),
			LinkPreviewOptions: &models.LinkPreviewOptions{
				IsDisabled: botBool(true),
			},
			ReplyMarkup: listingKeyboard(listing),
		})
		if err != nil {
			log.Printf("links: send listing id=%d: %v", listing.ID, err)
		}
	}
}

func (b *Bot) linksCallbackHandler(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.CallbackQuery == nil {
		return
	}

	cq := update.CallbackQuery
	chatID := cq.From.ID
	if banned, err := b.bans.IsBanned(ctx, chatID); err != nil {
		log.Printf("ban check chat_id=%d: %v", chatID, err)
	} else if banned {
		_, _ = b.api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: cq.ID,
			Text:            "Доступ закрыт",
			ShowAlert:       true,
		})
		return
	}
	action, listingID, ok := parseLinksCallback(cq.Data)
	if !ok {
		_, _ = b.api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: cq.ID,
			Text:            "Некорректная кнопка",
			ShowAlert:       true,
		})
		return
	}

	msg := cq.Message.Message
	if msg == nil {
		_, _ = b.api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: cq.ID,
			Text:            "Сообщение недоступно",
			ShowAlert:       true,
		})
		return
	}

	switch action {
	case "p":
		listing, err := b.listings.SetPaused(ctx, listingID, chatID, true)
		if !b.answerPause(ctx, cq.ID, err, listing, true) {
			return
		}
		b.refreshListingMessage(ctx, msg.Chat.ID, msg.ID, listing)

	case "r":
		listing, err := b.listings.SetPaused(ctx, listingID, chatID, false)
		if !b.answerPause(ctx, cq.ID, err, listing, false) {
			return
		}
		b.refreshListingMessage(ctx, msg.Chat.ID, msg.ID, listing)

	case "d":
		ok, err := b.listings.Delete(ctx, listingID, chatID)
		if err != nil || !ok {
			_, _ = b.api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: cq.ID,
				Text:            "Не удалось удалить",
				ShowAlert:       true,
			})
			return
		}
		_, _ = b.api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: cq.ID,
			Text:            "Удалено",
		})
		_, _ = b.api.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    msg.Chat.ID,
			MessageID: msg.ID,
			Text:      "Ссылка удалена.",
		})

	default:
		_, _ = b.api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: cq.ID,
			Text:            "Неизвестное действие",
			ShowAlert:       true,
		})
	}
}

func (b *Bot) answerPause(ctx context.Context, callbackID string, err error, listing *appmodels.Listing, pausing bool) bool {
	if errors.Is(err, repository.ErrBusy) {
		_, _ = b.api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callbackID,
			Text:            "Подожди секунду",
		})
		return false
	}
	fail := "Не удалось возобновить"
	ok := "Снова активна"
	if pausing {
		fail = "Не удалось поставить на паузу"
		ok = "На паузе"
	}
	if err != nil || listing == nil {
		_, _ = b.api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: callbackID,
			Text:            fail,
			ShowAlert:       true,
		})
		return false
	}
	_, _ = b.api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackID,
		Text:            ok,
	})
	return true
}

func (b *Bot) refreshListingMessage(ctx context.Context, chatID int64, messageID int, listing *appmodels.Listing) {
	_, err := b.api.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    chatID,
		MessageID: messageID,
		Text:      formatListingMessage(0, *listing),
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: botBool(true),
		},
		ReplyMarkup: listingKeyboard(*listing),
	})
	if err != nil {
		log.Printf("links: edit message chat_id=%d listing_id=%d: %v", chatID, listing.ID, err)
	}
}

func formatListingMessage(index int, listing appmodels.Listing) string {
	status := "активна"
	if listing.Paused {
		status = "пауза"
	}
	url := shortenURL(listing.URL)
	if index > 0 {
		return fmt.Sprintf("%d. [%s]\n%s", index, status, url)
	}
	return fmt.Sprintf("[%s]\n%s", status, url)
}

func listingKeyboard(listing appmodels.Listing) *models.InlineKeyboardMarkup {
	pauseText := "Пауза"
	pauseData := fmt.Sprintf("%sp:%d", linksCallbackPrefix, listing.ID)
	if listing.Paused {
		pauseText = "Возобновить"
		pauseData = fmt.Sprintf("%sr:%d", linksCallbackPrefix, listing.ID)
	}

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Открыть", URL: listing.URL},
				{Text: pauseText, CallbackData: pauseData},
				{Text: "Удалить", CallbackData: fmt.Sprintf("%sd:%d", linksCallbackPrefix, listing.ID)},
			},
		},
	}
}

func parseLinksCallback(data string) (action string, listingID int64, ok bool) {
	if !strings.HasPrefix(data, linksCallbackPrefix) {
		return "", 0, false
	}
	rest := strings.TrimPrefix(data, linksCallbackPrefix)
	parts := strings.Split(rest, ":")
	if len(parts) != 2 {
		return "", 0, false
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || id <= 0 {
		return "", 0, false
	}
	switch parts[0] {
	case "p", "r", "d":
		return parts[0], id, true
	default:
		return "", 0, false
	}
}

func shortenURL(raw string) string {
	const max = 80
	if len(raw) <= max {
		return raw
	}
	return raw[:max-1] + "…"
}

func (b *Bot) statusHandler(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	if b.rejectIfBanned(ctx, chatID) {
		return
	}
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
		"Твой тариф: %s\nПроверка объявлений: %s\nСсылок поиска: до %d",
		plan.Label(userPlan),
		plan.FormatIntervalRU(interval),
		plan.LinkLimit(userPlan),
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
	if b.rejectIfBanned(ctx, update.Message.Chat.ID) {
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

func (b *Bot) rejectIfBanned(ctx context.Context, chatID int64) bool {
	ban, err := b.bans.Get(ctx, chatID)
	if err != nil {
		log.Printf("ban check chat_id=%d: %v", chatID, err)
		return false
	}
	if ban == nil {
		return false
	}
	b.sendBanNotice(ctx, chatID, ban.Reason)
	return true
}

func (b *Bot) sendBanNotice(ctx context.Context, chatID int64, reason string) {
	text := "Доступ к FlatStalker сейчас закрыт."
	if reason = strings.TrimSpace(reason); reason != "" {
		text += "\nПричина: " + reason
	}
	text += fmt.Sprintf("\nЕсли это ошибка — напиши %s, разберёмся.", b.config.Telegram.SupportContact)
	_, _ = b.api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
}

func botBool(v bool) *bool {
	return new(v)
}
