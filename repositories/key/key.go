package key

import (
	"git.shi.foo/database"
	"git.shi.foo/models"
)

func Create(record *models.PublicKey) error {
	return database.DB.Create(record).Error
}

func ListByUser(userID uint) ([]models.PublicKey, error) {
	var records []models.PublicKey
	result := database.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&records)
	return records, result.Error
}

func FindByID(id uint) (*models.PublicKey, error) {
	var found models.PublicKey
	result := database.DB.First(&found, id)
	return &found, result.Error
}

func FindByFingerprint(fingerprint string) (*models.PublicKey, error) {
	var found models.PublicKey
	result := database.DB.Where("fingerprint = ?", fingerprint).First(&found)
	return &found, result.Error
}

func Update(record *models.PublicKey) error {
	return database.DB.Save(record).Error
}

func Delete(record *models.PublicKey) error {
	return database.DB.Unscoped().Delete(record).Error
}
