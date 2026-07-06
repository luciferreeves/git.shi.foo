package models

import "time"

type Cache struct {
	Key       string    `gorm:"primaryKey"`
	Value     []byte    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"index;not null"`
}
