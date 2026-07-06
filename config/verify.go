package config

import (
	"errors"
	"os"
	"path/filepath"
)

func verifyConfig() error {
	if GitHub.ClientID == "" {
		return errors.New(ClientIDRequired)
	}

	if GitHub.ClientSecret == "" {
		return errors.New(ClientSecretRequired)
	}

	if GitHub.EncryptionKey == "" {
		return errors.New(EncryptionKeyRequired)
	}

	if Git.ReposRoot == "" {
		homeDir, homeDirError := os.UserHomeDir()
		if homeDirError != nil {
			return homeDirError
		}
		Git.ReposRoot = filepath.Join(homeDir, DefaultReposDir)
	}

	return os.MkdirAll(Git.ReposRoot, DefaultDirPermission)
}
