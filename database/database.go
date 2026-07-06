package database

import (
	"git.shi.foo/config"
	"git.shi.foo/utils/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	dialector := postgres.Open(config.Database.DSN)

	var connectionError error
	DB, connectionError = gorm.Open(dialector, &gorm.Config{
		Logger: resolveGORMLogLevel(),
	})

	if connectionError != nil {
		logger.Fatalf(LogPrefix, ConnectionFailed, connectionError)
	}

	sqlDB, poolError := DB.DB()
	if poolError != nil {
		logger.Fatalf(LogPrefix, PoolConfigFailed, poolError)
	}

	sqlDB.SetMaxOpenConns(MaxOpenConnections)
	sqlDB.SetMaxIdleConns(MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(MaxConnectionLifetime)

	logger.Successf(LogPrefix, Connected, config.Database.DSN)

	migrate()
}
