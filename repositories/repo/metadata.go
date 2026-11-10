package repo

import (
	"git.shi.foo/database"
	"git.shi.foo/models"
)

func UpdateMetadata(record *models.Repo) error {
	return database.DB.Model(&models.Repo{}).
		Where("id = ?", record.ID).
		Select("Description", "Private", "DefaultBranch").
		Updates(record).Error
}
