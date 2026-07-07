package github

const (
	CiphertextTooShort      = "ciphertext too short"
	PrivateKeyInvalid       = "github app private key is not valid PEM"
	InstallationTokenFailed = "failed to mint installation token: status %d"
	UserKeysFetchFailed     = "failed to fetch user keys: status %d"
	UserEmailsFetchFailed   = "failed to fetch user emails: status %d"
	UserReposFetchFailed    = "failed to fetch user repos: status %d"
	RepoFetchFailed         = "failed to fetch repo: status %d"
)
