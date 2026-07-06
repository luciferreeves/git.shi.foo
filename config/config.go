package config

import (
	"git.shi.foo/utils/env"
	"git.shi.foo/utils/logger"
)

var (
	Server   server
	Database database
	Session  session
	GitHub   github
	Git      git
)

func init() {
	if err := env.Parse(&Server); err != nil {
		logger.Fatalf(LogPrefix, ServerConfigFailed, err)
	}

	if err := env.Parse(&Database); err != nil {
		logger.Fatalf(LogPrefix, DatabaseConfigFailed, err)
	}

	if err := env.Parse(&Session); err != nil {
		logger.Fatalf(LogPrefix, SessionConfigFailed, err)
	}

	if err := env.Parse(&GitHub); err != nil {
		logger.Fatalf(LogPrefix, GitHubConfigFailed, err)
	}

	if err := env.Parse(&Git); err != nil {
		logger.Fatalf(LogPrefix, GitConfigFailed, err)
	}

	if Server.Debug {
		logger.SetDebug(true)
	}

	if err := verifyConfig(); err != nil {
		logger.Fatalf(LogPrefix, VerificationFailed, err)
	}

	logger.Successf(LogPrefix, ConfigLoaded)
}
