package credential

import (
	"git.shi.foo/database"
	"git.shi.foo/models"
)

func RequireReauthorization(userID uint) error {
	return database.DB.Model(&models.Credential{}).
		Where("user_id = ?", userID).
		Update("reauthorization_required", true).Error
}
