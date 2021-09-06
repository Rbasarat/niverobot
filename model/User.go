package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int `gorm:"primaryKey"`
	Username  sql.NullString `gorm:"default=null;"`
	FirstName sql.NullString `gorm:"default=null;"`
	LastName sql.NullString `gorm:"default=null;"`
	LanguageCode sql.NullString `gorm:"default=null;"`
	IsBot	bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
