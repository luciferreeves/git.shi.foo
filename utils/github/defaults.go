package github

import "time"

const (
	LogPrefix       = "GitHub"
	StateLength     = 32
	StateAlphabet   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	StateKey        = "oauth_state"
	QueryParamCode  = "code"
	QueryParamState = "state"
	APIBase         = "https://api.github.com"
	UserEndpoint    = APIBase + "/user"
	PublicKeysURL   = "https://github.com/%s.keys"
	APIVersion      = "2022-11-28"
	AppJWTLifetime  = 9 * time.Minute
	AppJWTSkew      = 60 * time.Second
)

var Scopes = []string{"read:user", "user:email"}
