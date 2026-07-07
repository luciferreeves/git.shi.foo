package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Mirror(owner string, name string, token string) error {
	destination := RepoPath(owner, name)

	_ = os.RemoveAll(destination)

	if mkdirError := os.MkdirAll(filepath.Dir(destination), DirectoryPermission); mkdirError != nil {
		return fmt.Errorf(MkdirFailed, mkdirError)
	}

	cleanURL := fmt.Sprintf(CleanCloneURLFormat, owner, name)

	cloneCommand := exec.Command(
		"git",
		"-c", "credential.helper=",
		"-c", CredentialHelper,
		"clone", "--mirror", cleanURL, destination,
	)
	cloneCommand.Env = append(os.Environ(), GitTokenEnv+"="+token)

	if output, cloneError := cloneCommand.CombinedOutput(); cloneError != nil {
		scrubbed := strings.ReplaceAll(string(output), token, RedactedToken)
		return fmt.Errorf(CloneFailed, cloneError, scrubbed)
	}

	return nil
}
