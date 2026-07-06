package user

import (
	"git.shi.foo/database"
	"git.shi.foo/models"
)

func FindByProviderID(providerID string) (*models.User, error) {
	var found models.User
	result := database.DB.Where("provider_id = ?", providerID).First(&found)
	return &found, result.Error
}

func Create(record *models.User) error {
	return database.DB.Create(record).Error
}

func Update(record *models.User) error {
	return database.DB.Save(record).Error
}
