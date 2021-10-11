package model

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Kudos struct{}

type Kudo struct {
	ID         uint `gorm:"primaryKey"`
	IsPositive bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	MessageID  int   `gorm:"not null;"`
	ChatID     int64 `gorm:"not null"`
	UserID     int   `gorm:"not null;constraint:OnDelete:CASCADE"`
	User       User  `gorm:"not null;constraint:OnDelete:CASCADE"`
}

func (k Kudos) UpsertKudo(update *tgbotapi.Message, db *gorm.DB) (Kudo, bool, error) {
	var kudo Kudo
	result := db.Where(&Kudo{MessageID: update.ReplyToMessage.MessageID, ChatID: update.Chat.ID, UserID: update.From.ID}).Find(&kudo)

	if (strings.EqualFold(update.Text, "+") == kudo.IsPositive) && (strings.EqualFold(update.Text, "-") == !kudo.IsPositive) {
		return kudo, false, nil
	}

	isUpdate := false
	// Check if kudo does not exist on message and create
	if result.RowsAffected < 1 {
		kudo = Kudo{
			IsPositive: strings.EqualFold(update.Text, "+"),
			MessageID:  update.ReplyToMessage.MessageID,
			ChatID:     update.Chat.ID,
			UserID:     update.From.ID,
		}
		result = db.Create(&kudo)
	} else {
		db.Model(&kudo).Updates(map[string]interface{}{"is_positive": strings.EqualFold(update.Text, "+")})
		db.Where(&Kudo{MessageID: update.ReplyToMessage.MessageID, ChatID: update.Chat.ID, UserID: update.From.ID}).Find(&kudo)
		isUpdate = true
	}

	return kudo, isUpdate, result.Error
}
