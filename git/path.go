package git

import (
	"path/filepath"

	"git.shi.foo/config"
)

func RepoPath(owner string, name string) string {
	return filepath.Join(config.Git.ReposRoot, owner, name+RepoSuffix)
}
