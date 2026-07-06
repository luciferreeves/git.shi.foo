package config

const (
	ServerConfigFailed   = "Failed to parse server config: %v"
	DatabaseConfigFailed = "Failed to parse database config: %v"
	SessionConfigFailed  = "Failed to parse session config: %v"
	GitHubConfigFailed   = "Failed to parse github config: %v"
	GitConfigFailed      = "Failed to parse git config: %v"
	VerificationFailed   = "Configuration verification failed: %v"
	ConfigLoaded         = "Configuration loaded successfully."

	ClientIDRequired      = "GitHub Client ID is required"
	ClientSecretRequired  = "GitHub Client Secret is required"
	EncryptionKeyRequired = "GitHub encryption key is required"
)
