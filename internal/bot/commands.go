package bot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	appmodels "flat-stalker/internal/models"
	"flat-stalker/internal/plan"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const linksCallbackPrefix = "lnk:"

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
		if err != nil || listing == nil {
			_, _ = b.api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: cq.ID,
				Text:            "Не удалось поставить на паузу",
				ShowAlert:       true,
			})
			return
		}
		_, _ = b.api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: cq.ID,
			Text:            "На паузе",
		})
		b.refreshListingMessage(ctx, msg.Chat.ID, msg.ID, listing)

	case "r":
		listing, err := b.listings.SetPaused(ctx, listingID, chatID, false)
		if err != nil || listing == nil {
			_, _ = b.api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: cq.ID,
				Text:            "Не удалось возобновить",
				ShowAlert:       true,
			})
			return
		}
		_, _ = b.api.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: cq.ID,
			Text:            "Снова активна",
		})
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
