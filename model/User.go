package model

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
	"time"
)

type Users struct{}

type User struct {
	ID           int            `gorm:"primaryKey"`
	Username     sql.NullString `gorm:"default=null;"`
	FirstName    sql.NullString `gorm:"default=null;"`
	LastName     sql.NullString `gorm:"default=null;"`
	LanguageCode sql.NullString `gorm:"default=null;"`
	IsBot        bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u Users) CreateUserIfNotExist(update *tgbotapi.User, db *gorm.DB) (User, error) {
	var user User
	var result *gorm.DB
	if result = db.Take(&user, update.ID); result.RowsAffected < 1 {
		userNameNullable := sql.NullString{String: update.UserName, Valid: true}
		firstNameNullable := sql.NullString{String: update.FirstName, Valid: true}
		lastNameNullable := sql.NullString{String: update.LastName, Valid: true}
		LanguageCodeNullable := sql.NullString{String: update.LanguageCode, Valid: true}

		if update.UserName == "" {
			userNameNullable.Valid = false
		}
		if update.FirstName == "" {
			firstNameNullable.Valid = false
		}
		if update.LastName == "" {
			lastNameNullable.Valid = false
		}
		if update.LanguageCode == "" {
			LanguageCodeNullable.Valid = false
		}

		user = User{
			ID:           update.ID,
			Username:     userNameNullable,
			FirstName:    firstNameNullable,
			LastName:     lastNameNullable,
			LanguageCode: lastNameNullable,
			IsBot:        update.IsBot,
		}

		result = db.Create(&user)
	}

	return user, result.Error
}

func (u Users) Find(db *gorm.DB, id int) (User, error) {
	var user User
	result := db.Where(&User{ID: id}).Find(&user)
	return user, result.Error
}
