package database

import (
	"git.shi.foo/models"
	"git.shi.foo/utils/logger"
)

func migrate() {
	migrationError := DB.AutoMigrate(
		&models.User{},
		&models.Cache{},
	)

	if migrationError != nil {
		logger.Fatalf(LogPrefix, MigrationFailed, migrationError)
	}
}
