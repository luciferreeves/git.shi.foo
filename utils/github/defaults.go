package github

import "time"

const (
	LogPrefix          = "GitHub"
	StateLength        = 32
	StateAlphabet      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	StateKey           = "oauth_state"
	QueryParamCode     = "code"
	QueryParamState    = "state"
	APIBase            = "https://api.github.com"
	UserEndpoint       = APIBase + "/user"
	UserKeysEndpoint   = APIBase + "/user/keys"
	UserEmailsEndpoint = APIBase + "/user/emails"
	UserReposEndpoint  = APIBase + "/user/repos?per_page=100&sort=updated"
	RepoEndpointFormat = APIBase + "/repos/%s/%s"
	APIVersion         = "2022-11-28"
	AppJWTLifetime     = 9 * time.Minute
	AppJWTSkew         = 60 * time.Second
)

var Scopes = []string{"read:user", "user:email"}
