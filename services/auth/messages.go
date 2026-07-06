package auth

const (
	StateGenerationLog = "Failed to generate OAuth state: %v"
	TokenExchangeLog   = "Failed to exchange OAuth code: %v"
	IdentityFetchLog   = "Failed to fetch GitHub identity: %v"
	UserUpsertLog      = "Failed to persist user: %v"
	CredentialStoreLog = "Failed to store refresh token: %v"

	StateGenerationFailed = "Could not start GitHub sign-in."
	TokenExchangeFailed   = "GitHub sign-in failed."
	IdentityFetchFailed   = "Could not read your GitHub identity."
	UserUpsertFailed      = "Could not save your account."
	CredentialStoreFailed = "Could not store your credentials."
)
