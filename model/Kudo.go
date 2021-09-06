package model

import "time"

type Kudo struct {
	ID         uint `gorm:"primaryKey"`
	IsPositive bool
	CreatedAt  time.Time
	MessageID  int   `gorm:"not null;autoIncrement:false"`
	ChatID     int64 `gorm:"not null"`
	UserID     int  `gorm:"not null;constraint:OnDelete:CASCADE"`
	User       User  `gorm:"not null;constraint:OnDelete:CASCADE"`
}
