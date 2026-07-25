package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *Bot) startHandler(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	user, err := b.users.CreateByChatID(ctx, chatID)
	if err != nil {
		log.Printf("start: create user chat_id=%d: %v", chatID, err)
		_, _ = b.api.SendMessage(b.ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Не удалось сохранить профиль. Попробуй ещё раз чуть позже.",
		})
		return
	}

	log.Printf("start: user id=%d chat_id=%d", user.ID, user.ChatID)

	_, err = b.api.SendMessage(b.ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Привет! Ты в FlatStalker. Открой кабинет через кнопку меню и настрой фильтры.",
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
		Text:      fmt.Sprintf("По всем вопросам пишите сюда <b>%s</b>", b.config.Telegram.SupportContact),
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		log.Println(err)
	}
}
