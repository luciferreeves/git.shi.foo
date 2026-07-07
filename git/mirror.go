package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Mirror(owner string, name string, token string) error {
	destination := RepoPath(owner, name)

	_ = os.RemoveAll(destination)

	if mkdirError := os.MkdirAll(filepath.Dir(destination), DirectoryPermission); mkdirError != nil {
		return fmt.Errorf(MkdirFailed, mkdirError)
	}

	authURL := fmt.Sprintf(AuthCloneURLFormat, token, owner, name)
	cloneCommand := exec.Command("git", "clone", "--mirror", authURL, destination)
	if output, cloneError := cloneCommand.CombinedOutput(); cloneError != nil {
		return fmt.Errorf(CloneFailed, cloneError, string(output))
	}

	cleanURL := fmt.Sprintf(CleanCloneURLFormat, owner, name)
	resetCommand := exec.Command("git", "-C", destination, "remote", "set-url", "origin", cleanURL)
	_ = resetCommand.Run()

	return nil
}
