package settings

const (
	TokenGenerateLog = "Failed to generate token: %v"
	TokenCreateLog   = "Failed to create token: %v"
	TokenRevokeLog   = "Failed to revoke token: %v"
	TokenListLog     = "Failed to list tokens: %v"
	KeyAddLog        = "Failed to add public key: %v"
	KeyRemoveLog     = "Failed to remove public key: %v"
	KeyListLog       = "Failed to list public keys: %v"
	KeyImportLog     = "Failed to fetch GitHub public keys: %v"
	KeyImportSkipLog = "Skipped an unparseable imported key: %s"

	AuthRequired      = "Sign in required."
	LabelRequired     = "A label is required."
	TokenCreateFailed = "Could not create the token."
	TokenRevokeFailed = "Could not revoke the token."
	TokenListFailed   = "Could not load your tokens."
	TokenNotFound     = "Token not found."
	TokenNotYours     = "That token is not yours."

	InvalidPublicKey = "That doesn't look like a valid SSH public key."
	KeyExists        = "That key is already registered."
	KeyAddFailed     = "Could not add the key."
	KeyListFailed    = "Could not load your keys."
	KeyNotFound      = "Key not found."
	KeyNotYours      = "That key is not yours."
	KeyRemoveFailed  = "Could not remove the key."
	KeyImportFailed  = "Could not fetch your keys from GitHub."
)
