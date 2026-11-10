package credentials

import "time"

const (
	LogPrefix                            = "Credentials"
	AccessCacheKeyFormat                 = "github:access:%d"
	AccessSafetyMargin                   = 60 * time.Second
	AccessFallbackTTL                    = 8 * time.Hour
	InvalidRefreshTokenCode              = "bad_refresh_token"
	InvalidGrantCode                     = "invalid_grant"
	InstallationCacheKeyFormat           = "github:installation:%s"
	RepositoryInstallationCacheKeyFormat = "github:repository-installation:%s/%s"
	RepositoryInstallationTTL            = time.Hour
)
