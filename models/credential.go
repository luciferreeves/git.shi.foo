package models

import "gorm.io/gorm"

type Credential struct {
	gorm.Model
	UserID       uint   `gorm:"uniqueIndex;not null"`
	RefreshToken string `gorm:"type:text;not null"`
}
