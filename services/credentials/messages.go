package credentials

const (
	CredentialMissingLog    = "No stored credential for user %d: %v."
	DecryptLog              = "Failed to decrypt refresh token for user %d: %v."
	RefreshLog              = "Failed to refresh GitHub token for user %d: %v."
	RotateLog               = "Failed to persist refreshed GitHub tokens for user %d: %v."
	CacheLog                = "Failed to clear access-token cache for user %d: %v."
	ReauthorizationRequired = "This GitHub account must log in again to reconnect."
	ReauthorizationLog      = "GitHub rejected the refresh token for user %d. Log in with that GitHub account to reconnect."
	ReauthorizationStoreLog = "Failed to mark GitHub user %d for reauthorization: %v."
)
