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

func Count() (int64, error) {
	var total int64
	result := database.DB.Model(&models.User{}).Count(&total)
	return total, result.Error
}

func All() ([]models.User, error) {
	var records []models.User
	result := database.DB.Order("created_at desc").Find(&records)
	return records, result.Error
}
