package models

import (
	"time"

	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	Kind        string `gorm:"index;not null"`
	RepoID      *uint  `gorm:"index"`
	Status      string `gorm:"index;not null;default:pending"`
	Phase       string
	Percent     int `gorm:"not null;default:0"`
	Attempts    int `gorm:"not null;default:0"`
	MaxAttempts int `gorm:"not null;default:3"`
	RunAfter    time.Time
	Error       string `gorm:"type:text"`
}
