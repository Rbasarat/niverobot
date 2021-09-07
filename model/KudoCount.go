package model

import (
	"gorm.io/gorm"
	"time"
)

type KudoCounts struct{}

type KudoCount struct {
	UserID    int   `gorm:"not null;constraint:OnDelete:CASCADE;primaryKey"`
	ChatID    int64 `gorm:"primaryKey"`
	Plus      int   `gorm:"default:0"`
	Minus     int   `gorm:"default:0"`
	UpdatedAt time.Time
	User      User `gorm:"not null;constraint:OnDelete:CASCADE"`
}
func (k KudoCounts) GetKudoCountPerChat(chatId int64, db *gorm.DB) ([]KudoCount, error) {
	var kudoCounts []KudoCount
	result := db.Where(&KudoCount{ChatID: chatId}).Preload("User").Order("plus desc").Find(&kudoCounts)
	return kudoCounts, result.Error
}

func (k KudoCounts) GetKudoCount(userId int, chatId int64, db *gorm.DB) (KudoCount, error) {
	var kudoCount KudoCount
	result := db.Where(&KudoCount{UserID: userId, ChatID: chatId}).Find(&kudoCount)
	return kudoCount, result.Error
}

func (k KudoCounts) UpdateKudoCount(kudo Kudo, user User, db *gorm.DB, chatId int64) (KudoCount, error) {
	var kudoCount KudoCount
	result := db.Where(&KudoCount{UserID: user.ID}).Find(&kudoCount)

	if result.RowsAffected < 1 {
		kudoCount = KudoCount{
			User:   user,
			ChatID: chatId,
			Plus:   0,
			Minus:  0,
		}
	}

	if kudo.IsPositive {
		kudoCount.Plus++
	} else {
		kudoCount.Minus++
	}

	result = db.Save(&kudoCount)
	return kudoCount, result.Error

}
