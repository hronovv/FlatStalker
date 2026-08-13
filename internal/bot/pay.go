package bot

import (
	"context"
	"fmt"
	"log"

	"flat-stalker/internal/plan"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (b *Bot) BotURL() string {
	if b.username == "" {
		return ""
	}
	return "https://t.me/" + b.username
}

func (b *Bot) IssuePlanInvoice(ctx context.Context, chatID int64, payload, title, description string, amountKop int) (string, error) {
	_, err := b.api.SendInvoice(ctx, &bot.SendInvoiceParams{
		ChatID:        chatID,
		Title:         title,
		Description:   description,
		Payload:       payload,
		ProviderToken: b.config.Telegram.PaymentProviderToken,
		Currency:      plan.Currency,
		Prices: []models.LabeledPrice{
			{Label: plan.FormatBYN(amountKop) + " BYN", Amount: amountKop},
		},
	})
	if err != nil {
		return "", err
	}
	return b.BotURL(), nil
}

func (b *Bot) preCheckoutHandler(ctx context.Context, _ *bot.Bot, update *models.Update) {
	q := update.PreCheckoutQuery
	if q == nil {
		return
	}

	ok := true
	errMsg := ""
	payment, err := b.payments.GetByPayload(ctx, q.InvoicePayload)
	if err != nil {
		log.Printf("precheckout: payload=%s: %v", q.InvoicePayload, err)
		ok = false
		errMsg = "Не удалось проверить оплату. Попробуй ещё раз."
	} else if payment == nil || payment.Status != "pending" {
		ok = false
		errMsg = "Счёт уже недействителен. Создай оплату заново."
	} else if q.TotalAmount != payment.AmountKop {
		ok = false
		errMsg = "Сумма не совпадает. Создай оплату заново."
	} else if q.From != nil && payment.ChatID != 0 && q.From.ID != payment.ChatID {
		ok = false
		errMsg = "Это оплата другого аккаунта."
	}

	if _, err := b.api.AnswerPreCheckoutQuery(ctx, &bot.AnswerPreCheckoutQueryParams{
		PreCheckoutQueryID: q.ID,
		OK:                 ok,
		ErrorMessage:       errMsg,
	}); err != nil {
		log.Printf("precheckout answer: %v", err)
	}
}

func (b *Bot) successfulPaymentHandler(ctx context.Context, _ *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.SuccessfulPayment == nil {
		return
	}

	sp := update.Message.SuccessfulPayment
	chatID := update.Message.Chat.ID
	payment, err := b.payments.MarkPaid(ctx, sp.InvoicePayload, sp.TelegramPaymentChargeID)
	if err != nil {
		log.Printf("payment: mark paid payload=%s: %v", sp.InvoicePayload, err)
		return
	}
	if payment == nil {
		return
	}

	if err := b.users.ApplyPlan(ctx, chatID, payment.Plan, payment.Days); err != nil {
		log.Printf("payment: apply plan chat_id=%d: %v", chatID, err)
		return
	}

	text := fmt.Sprintf(
		"Оплата прошла. Тариф %s на %s включён.",
		plan.Label(payment.Plan),
		plan.FormatDaysRU(payment.Days),
	)
	if _, err := b.api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	}); err != nil {
		log.Printf("payment: notify chat_id=%d: %v", chatID, err)
	}
}
