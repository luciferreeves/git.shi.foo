package github

import (
	"git.shi.foo/config"

	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

var OAuthConfig oauth2.Config

func init() {
	OAuthConfig = oauth2.Config{
		ClientID:     config.GitHub.ClientID,
		ClientSecret: config.GitHub.ClientSecret,
		RedirectURL:  config.GitHub.CallbackURL,
		Endpoint:     githuboauth.Endpoint,
		Scopes:       Scopes,
	}
}
