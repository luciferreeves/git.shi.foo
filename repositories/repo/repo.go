package repo

import (
	"git.shi.foo/database"
	"git.shi.foo/models"
)

func Create(record *models.Repo) error {
	return database.DB.Create(record).Error
}

func Update(record *models.Repo) error {
	return database.DB.Save(record).Error
}

func FindByID(id uint) (*models.Repo, error) {
	var found models.Repo
	result := database.DB.First(&found, id)
	return &found, result.Error
}

func FindByOwnerName(owner string, name string) (*models.Repo, error) {
	var found models.Repo
	result := database.DB.Where("owner = ? AND name = ?", owner, name).First(&found)
	return &found, result.Error
}

func ListActive() ([]models.Repo, error) {
	var records []models.Repo
	result := database.DB.Where("status = ?", StatusActive).Order("updated_at desc").Find(&records)
	return records, result.Error
}

func ListAll() ([]models.Repo, error) {
	var records []models.Repo
	result := database.DB.Order("updated_at desc").Find(&records)
	return records, result.Error
}
