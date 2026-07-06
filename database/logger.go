package database

import (
	"git.shi.foo/config"

	"gorm.io/gorm/logger"
)

func resolveGORMLogLevel() logger.Interface {
	switch config.Server.Debug {
	case true:
		return logger.Default.LogMode(logger.Info)
	default:
		return logger.Default.LogMode(logger.Silent)
	}
}
