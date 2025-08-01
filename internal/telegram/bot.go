package telegram

import (
	"fmt"

	"github.com/PhosFactum/budget-guardian/internal/interfaces"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	service interfaces.BudgetService
}

// NewBot: creating new bot
func NewBot(token string, service interfaces.BudgetService) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	bot.Debug = true
	return &Bot{api: bot, service: service}, nil
}

// Start: function to start a bot
func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		msg := tgbotapi.NewMessage(chatID, "")

		switch update.Message.Command() {
		case "start":
			msg.Text = "Добро пожаловать в Budget Guardian!\n\n" +
				"Доступные команды:\n" +
				"/setup - Настройка бюджета\n" +
				"/spent [сумма] - Записать трату\n" +
				"/limit - Показать дневной лимит"
		case "setup":
			msg.Text = b.handleSetup(update.Message)
		case "spent":
			msg.Text = b.handleSpent(update.Message)
		case "limit":
			msg.Text = b.handleLimit(update.Message)
		default:
			msg.Text = "Неизвестная команда"
		}

		if _, err := b.api.Send(msg); err != nil {
			fmt.Printf("Failed to send message: %v\n", err)
		}
	}
}

func (b *Bot) handleSetup(msg *tgbotapi.Message) string {
	// Configuration logic
	_ = msg
	return "Настройки сохранены! Пример:\n" +
		"Доход: 20000₽\n" +
		"Фикс.расходы: 5000₽\n" +
		"Дни выплат: 5 и 20 числа"
}

func (b *Bot) handleSpent(msg *tgbotapi.Message) string {
	// Transaction recording logic
	_ = msg
	return "Трата записана! Текущий лимит: 750₽/день"
}

func (b *Bot) handleLimit(msg *tgbotapi.Message) string {
	_ = msg
	// Showint limit logic
	return "Ваш дневной лимит: 750₽\n" +
		"Остаток на сегодня: 320₽\n" +
		"Накопленный долг: -120₽"
}
