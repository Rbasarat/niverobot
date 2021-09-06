package features

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
)

func AddKudo(update *tgbotapi.Message, db *gorm.DB) error {

	user, err := CreateUserIfNotExist(update.From, db)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	kudo, err := CreateKudoIfNotExist(update, user, db)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	_, err = UpdateKudoCount(kudo, user, db)
	fmt.Printf("last step error: %s", err)

	return err
}
