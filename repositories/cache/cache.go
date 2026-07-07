package cache

import (
	"time"

	"git.shi.foo/database"
	"git.shi.foo/models"

	"gorm.io/gorm/clause"
)

func Get(key string) (*models.Cache, error) {
	var entry models.Cache
	result := database.DB.Where("key = ? AND expires_at > ?", key, time.Now()).First(&entry)
	return &entry, result.Error
}

func Set(key string, value []byte, expiresAt time.Time) error {
	entry := models.Cache{Key: key, Value: value, ExpiresAt: expiresAt}
	return database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "expires_at"}),
	}).Create(&entry).Error
}

func Delete(key string) error {
	return database.DB.Where("key = ?", key).Delete(&models.Cache{}).Error
}

func Sweep() error {
	return database.DB.Where("expires_at <= ?", time.Now()).Delete(&models.Cache{}).Error
}
