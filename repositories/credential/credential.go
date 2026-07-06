package credential

import (
	"errors"

	"git.shi.foo/database"
	"git.shi.foo/models"

	"gorm.io/gorm"
)

func FindByUserID(userID uint) (*models.Credential, error) {
	var found models.Credential
	result := database.DB.Where("user_id = ?", userID).First(&found)
	return &found, result.Error
}

func Upsert(userID uint, refreshToken string) error {
	existing, lookupError := FindByUserID(userID)
	if lookupError != nil {
		if errors.Is(lookupError, gorm.ErrRecordNotFound) {
			return database.DB.Create(&models.Credential{UserID: userID, RefreshToken: refreshToken}).Error
		}
		return lookupError
	}

	existing.RefreshToken = refreshToken
	return database.DB.Save(existing).Error
}
