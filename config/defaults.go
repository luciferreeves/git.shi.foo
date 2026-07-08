package config

import "os"

const (
	LogPrefix            = "Config"
	DefaultDirPermission = os.FileMode(0o755)
	DefaultReposDir      = ".git.shi.foo/repos"
	AppVersion           = "0.1.0"
)
