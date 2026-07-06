package models

import (
	"errors"
	"strings"

	"git.shi.foo/account"
	"git.shi.foo/utils/validate"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ProviderID string `gorm:"uniqueIndex;not null"`
	Login      string `gorm:"uniqueIndex;not null"`
	Email      string `gorm:"index"`
	Avatar     string `gorm:"type:text"`
	Admin      bool   `gorm:"not null;default:false"`
	Enabled    bool   `gorm:"not null;default:true"`
}

func (self *User) BeforeCreate(tx *gorm.DB) error {
	return self.normalize()
}

func (self *User) BeforeUpdate(tx *gorm.DB) error {
	return self.normalize()
}

func (self *User) normalize() error {
	self.Login = strings.TrimSpace(self.Login)
	self.Email = strings.TrimSpace(strings.ToLower(self.Email))

	if self.Email != "" && !validate.Email(self.Email) {
		return errors.New(validate.InvalidEmail)
	}

	return nil
}

func (self *User) ToResponse() account.Response {
	return account.Response{
		ID:         self.ID,
		ProviderID: self.ProviderID,
		Login:      self.Login,
		Email:      self.Email,
		Avatar:     self.Avatar,
		Admin:      self.Admin,
		Enabled:    self.Enabled,
		CreatedAt:  self.CreatedAt,
	}
}
