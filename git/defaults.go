package git

import "os"

const (
	LogPrefix           = "Git"
	DirectoryPermission = os.FileMode(0o755)
	RepoSuffix          = ".git"
	AuthCloneURLFormat  = "https://x-access-token:%s@github.com/%s/%s.git"
	CleanCloneURLFormat = "https://github.com/%s/%s.git"
)
