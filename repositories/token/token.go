package token

import (
	"git.shi.foo/database"
	"git.shi.foo/models"
)

func Create(record *models.PersonalAccessToken) error {
	return database.DB.Create(record).Error
}

func ListByUser(userID uint) ([]models.PersonalAccessToken, error) {
	var records []models.PersonalAccessToken
	result := database.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&records)
	return records, result.Error
}

func FindByID(id uint) (*models.PersonalAccessToken, error) {
	var found models.PersonalAccessToken
	result := database.DB.First(&found, id)
	return &found, result.Error
}

func FindByHash(hash string) (*models.PersonalAccessToken, error) {
	var found models.PersonalAccessToken
	result := database.DB.Where("token_hash = ?", hash).First(&found)
	return &found, result.Error
}

func Delete(record *models.PersonalAccessToken) error {
	return database.DB.Delete(record).Error
}
