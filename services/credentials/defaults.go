package credentials

import "time"

const (
	LogPrefix            = "Credentials"
	AccessCacheKeyFormat = "github:access:%d"
	AccessSafetyMargin   = 60 * time.Second
	AccessFallbackTTL    = 8 * time.Hour
)
