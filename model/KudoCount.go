package model

import (
	"time"
)

type KudoCount struct {
	UserID    int `gorm:"not null;constraint:OnDelete:CASCADE;primaryKey"`
	Plus      int  `gorm:"default:0"`
	Minus     int  `gorm:"default:0"`
	UpdatedAt time.Time
	User      User `gorm:"not null;constraint:OnDelete:CASCADE"`
}
