package features

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
)

type Action interface {
	Execute(update tgbotapi.Update, db *gorm.DB, bot *tgbotapi.BotAPI)
	Trigger(update tgbotapi.Update) bool
}
