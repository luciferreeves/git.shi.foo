package models

import (
	"time"

	"gorm.io/gorm"
)

type PersonalAccessToken struct {
	gorm.Model
	UserID     uint   `gorm:"index;not null"`
	Label      string `gorm:"not null"`
	TokenHash  string `gorm:"uniqueIndex;not null"`
	Preview    string `gorm:"not null"`
	LastUsedAt *time.Time
}
