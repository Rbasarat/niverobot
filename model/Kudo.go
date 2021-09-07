package model

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Kudos struct {}

type Kudo struct {
	ID         uint `gorm:"primaryKey"`
	IsPositive bool
	CreatedAt  time.Time
	MessageID  int   `gorm:"not null;"`
	ChatID     int64 `gorm:"not null"`
	UserID     int  `gorm:"not null;constraint:OnDelete:CASCADE"`
	User       User  `gorm:"not null;constraint:OnDelete:CASCADE"`
}

func (k Kudos) CreateKudoIfNotExist(update *tgbotapi.Message, db *gorm.DB) (Kudo, error) {
	var kudo Kudo
	var result *gorm.DB
	// Check if kudo does not exist on message and create
	if result = db.Where(&Kudo{MessageID: update.ReplyToMessage.MessageID, ChatID: update.Chat.ID, UserID: update.From.ID}).Find(&kudo); result.RowsAffected < 1 {
		kudo = Kudo{
			IsPositive: strings.EqualFold(update.Text, "+"),
			MessageID:  update.ReplyToMessage.MessageID,
			ChatID:     update.Chat.ID,
			UserID:     update.From.ID,
		}
		result = db.Create(&kudo)
	} else {
		return kudo, errors.New("kudo already added")
	}

	return kudo, result.Error
}