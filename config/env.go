package config

import "time"

type server struct {
	Host  string `env:"HOST" default:"0.0.0.0"`
	Port  int    `env:"PORT" default:"3333"`
	Debug bool   `env:"DEBUG" default:"false"`
}

type database struct {
	DSN string `env:"DSN" default:"host=localhost user=postgres password=postgres dbname=gitshifoo port=5432 sslmode=disable"`
}

type session struct {
	CookieDomain   string        `env:"SESSION_COOKIE_DOMAIN" default:"localhost"`
	CookieName     string        `env:"SESSION_COOKIE_NAME" default:"gitshifoo_session"`
	CookiePath     string        `env:"SESSION_COOKIE_PATH" default:"/"`
	CookieSameSite string        `env:"SESSION_COOKIE_SAME_SITE" default:"Lax"`
	CookieSecure   bool          `env:"SESSION_SECURE_COOKIE" default:"false"`
	CookieTimeout  time.Duration `env:"SESSION_TIMEOUT" default:"24h"`
}

type github struct {
	ClientID       string `env:"GITHUB_CLIENT_ID" default:""`
	ClientSecret   string `env:"GITHUB_CLIENT_SECRET" default:""`
	CallbackURL    string `env:"GITHUB_CALLBACK_URL" default:"http://localhost:3333/auth/callback"`
	AppID          string `env:"GITHUB_APP_ID" default:""`
	AppPrivateKey  string `env:"GITHUB_APP_PRIVATE_KEY" default:""`
	InstallationID string `env:"GITHUB_INSTALLATION_ID" default:""`
	EncryptionKey  string `env:"GITHUB_ENCRYPTION_KEY" default:""`
}

type git struct {
	ReposRoot  string `env:"REPOS_ROOT" default:""`
	SSHHostKey string `env:"SSH_HOST_KEY" default:""`
	SSHAddress string `env:"SSH_ADDRESS" default:":2222"`
}
