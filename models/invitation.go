package models

import "gorm.io/gorm"

type Invitation struct {
	gorm.Model
	Email     string `gorm:"index;not null"`
	Username  string `gorm:"index;not null"`
	Token     string `gorm:"uniqueIndex;not null"`
	Status    string `gorm:"index;not null;default:pending"`
	InvitedBy uint   `gorm:"not null"`
}
