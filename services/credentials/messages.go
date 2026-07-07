package credentials

const (
	CredentialMissingLog = "No stored credential for user: %v"
	DecryptLog           = "Failed to decrypt refresh token: %v"
	RefreshLog           = "Failed to refresh GitHub token: %v"
	RotateLog            = "Failed to rotate refresh token: %v"
	CacheLog             = "Failed to cache access token: %v"
)
