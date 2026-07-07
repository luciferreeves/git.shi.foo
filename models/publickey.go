package models

import "gorm.io/gorm"

type PublicKey struct {
	gorm.Model
	UserID      uint   `gorm:"index;not null"`
	Title       string `gorm:"not null"`
	KeyType     string `gorm:"not null"`
	Fingerprint string `gorm:"uniqueIndex;not null"`
	Body        string `gorm:"type:text;not null"`
	Source      string `gorm:"not null"`
	GitHubID    int64  `gorm:"index"`
}
