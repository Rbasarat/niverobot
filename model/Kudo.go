package model

import (
	"time"
)

type Kudo struct {
	ID        uint `gorm:"primaryKey"`
	Plus      int  `gorm:"default:0"`
	Minus     int  `gorm:"default:0"`
	UpdatedAt time.Time
	UserID    uint `gorm:"not null;constraint:OnDelete:CASCADE"`
	User      User `gorm:"not null;constraint:OnDelete:CASCADE"`
}
