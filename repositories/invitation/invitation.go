package invitation

import (
	"git.shi.foo/database"
	"git.shi.foo/models"
)

func Create(record *models.Invitation) error {
	return database.DB.Create(record).Error
}

func Update(record *models.Invitation) error {
	return database.DB.Save(record).Error
}

func FindByToken(token string) (*models.Invitation, error) {
	var found models.Invitation
	result := database.DB.Where("token = ?", token).First(&found)
	return &found, result.Error
}

func FindByUsernameAndStatus(username string, status string) (*models.Invitation, error) {
	var found models.Invitation
	result := database.DB.Where("username = ? AND status = ?", username, status).First(&found)
	return &found, result.Error
}

func All() ([]models.Invitation, error) {
	var records []models.Invitation
	result := database.DB.Order("created_at desc").Find(&records)
	return records, result.Error
}
