package features

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
)

func AddKudo(update *tgbotapi.Message, db *gorm.DB) error {

	receiver, err := CreateUserIfNotExist(update.ReplyToMessage.From, db)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	_, err = CreateUserIfNotExist(update.From, db)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	kudo, err := CreateKudoIfNotExist(update, receiver, db)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	_, err = UpdateKudoCount(kudo, receiver, db)

	return err
}
