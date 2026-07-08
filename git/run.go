package git

import (
	"bytes"
	"fmt"
	"os/exec"
)

func runGit(repoPath string, arguments ...string) ([]byte, error) {
	command := exec.Command("git", append([]string{"-C", repoPath}, arguments...)...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	if runError := command.Run(); runError != nil {
		return nil, fmt.Errorf(GitCommandFailed, runError, stderr.String())
	}

	return stdout.Bytes(), nil
}
