package model

import (
	"time"
)

type KudoCount struct {
	Plus      int  `gorm:"default:0"`
	Minus     int  `gorm:"default:0"`
	UpdatedAt time.Time
	UserID    uint `gorm:"not null;constraint:OnDelete:CASCADE;primaryKey"`
	User      User `gorm:"not null;constraint:OnDelete:CASCADE"`
}
