package database

import (
	"git.shi.foo/models"
	"git.shi.foo/utils/logger"
)

func migrate() {
	migrationError := DB.AutoMigrate(
		&models.User{},
		&models.Cache{},
		&models.Credential{},
		&models.Invitation{},
		&models.PersonalAccessToken{},
		&models.PublicKey{},
		&models.Repo{},
	)

	if migrationError != nil {
		logger.Fatalf(LogPrefix, MigrationFailed, migrationError)
	}
}
