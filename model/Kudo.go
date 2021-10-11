package model

import (
	"errors"
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

func (k Kudos) UpsertKudo(text string, messageId int, userId int, chatId int64, db *gorm.DB) (Kudo, bool, error) {
	var kudo Kudo
	result := db.Where(&Kudo{MessageID: messageId, ChatID: chatId, UserID: userId}).Find(&kudo)

	if (strings.EqualFold(text, "+") == kudo.IsPositive) && (strings.EqualFold(text, "-") == !kudo.IsPositive) {
		return kudo, false, errors.New("kudo already exist")
	}

	isUpdate := false
	// Check if kudo does not exist on message and create
	if result.RowsAffected < 1 {
		kudo = Kudo{
			IsPositive: strings.EqualFold(text, "+"),
			MessageID:  messageId,
			ChatID:     chatId,
			UserID:     userId,
		}
		result = db.Create(&kudo)
	} else {
		db.Model(&kudo).Updates(map[string]interface{}{"is_positive": strings.EqualFold(text, "+")})
		db.Where(&Kudo{MessageID: messageId, ChatID: chatId, UserID: userId}).Find(&kudo)
		isUpdate = true
	}

	return kudo, isUpdate, result.Error
}
