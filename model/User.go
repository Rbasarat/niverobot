package model

import "time"

type User struct {
	ID        int64 `gorm:"primaryKey"`
	Username  string
	FirstName string
	CreatedAt time.Time
	UpdatedAt time.Time
}
