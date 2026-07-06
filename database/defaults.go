package database

import "time"

const (
	LogPrefix             = "Database"
	MaxConnectionLifetime = time.Hour
	MaxIdleConnections    = 5
	MaxOpenConnections    = 25
)
