package features

import (
	"database/sql"
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
	"niverobot/model"
	"strings"
)

func CreateUserIfNotExist(from *tgbotapi.User, db *gorm.DB) (model.User, error) {
	var user model.User
	var result *gorm.DB
	if result = db.Take(&user, from.ID); result.Error == gorm.ErrRecordNotFound {
		userNameNullable := sql.NullString{String: from.UserName, Valid: true}
		firstNameNullable := sql.NullString{String: from.FirstName, Valid: true}
		lastNameNullable := sql.NullString{String: from.LastName, Valid: true}
		LanguageCodeNullable := sql.NullString{String: from.LanguageCode, Valid: true}

		if from.UserName == "" {
			userNameNullable.Valid = false
		}
		if from.FirstName == "" {
			firstNameNullable.Valid = false
		}
		if from.LastName == "" {
			lastNameNullable.Valid = false
		}
		if from.LanguageCode == "" {
			LanguageCodeNullable.Valid = false
		}

		user = model.User{
			ID:           from.ID,
			Username:     userNameNullable,
			FirstName:    firstNameNullable,
			LastName:     lastNameNullable,
			LanguageCode: lastNameNullable,
			IsBot:        from.IsBot,
		}

		result = db.Create(&user)
	}

	return user, result.Error
}

func CreateKudoIfNotExist(update *tgbotapi.Message, user model.User, db *gorm.DB) (model.Kudo, error) {
	var kudo model.Kudo
	var result *gorm.DB
	kudoPositive := strings.EqualFold(update.Text, "+")
	// Check if kudo does not exist on message and create
	if result = db.Where(&model.Kudo{MessageID: update.ReplyToMessage.MessageID, ChatID: update.Chat.ID, User: user}).Find(&kudo); result.RowsAffected < 1 {
		kudo = model.Kudo{
			IsPositive: kudoPositive,
			MessageID:  update.ReplyToMessage.MessageID,
			ChatID:     update.Chat.ID,
			User:       user,
		}
		result = db.Create(&kudo)
	} else {
		return kudo, errors.New("kudo already added")
	}

	return kudo, result.Error
}

func UpdateKudoCount(kudo model.Kudo, user model.User, db *gorm.DB) (model.KudoCount, error) {
	var kudoCount model.KudoCount
	var result *gorm.DB

	// Update kudo counter
	if result = db.Where(&model.KudoCount{User: user}).Find(&kudoCount); result.Error == gorm.ErrRecordNotFound {
		kudoCount = model.KudoCount{
			User: user,
		}
		result = db.Create(&kudoCount)
	}
	if result.Error != nil {
		return kudoCount, result.Error
	}

	kudoCount.User = user
	if kudo.IsPositive {
		kudoCount.Plus++
	} else {
		kudoCount.Minus++
	}
	result = db.Save(&kudoCount)
	return kudoCount, result.Error

}
