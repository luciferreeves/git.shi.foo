package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Fetch(owner string, name string, token string, onProgress func(phase string, percent int)) error {
	command := exec.Command(
		"git",
		"-c", "credential.helper=",
		"-c", CredentialHelper,
		"-C", RepoPath(owner, name),
		"fetch", "--prune", "--progress", "origin",
	)
	command.Env = append(os.Environ(), GitTokenEnv+"="+token)

	stderrPipe, pipeError := command.StderrPipe()
	if pipeError != nil {
		return fmt.Errorf(FetchFailed, pipeError, "")
	}

	if startError := command.Start(); startError != nil {
		return fmt.Errorf(FetchFailed, startError, "")
	}

	var captured bytes.Buffer
	scanner := bufio.NewScanner(stderrPipe)
	scanner.Split(scanProgressLines)
	for scanner.Scan() {
		line := scanner.Text()
		captured.WriteString(line)
		captured.WriteByte('\n')

		if onProgress == nil {
			continue
		}
		if phase, percent, matched := parseProgress(line); matched {
			onProgress(phase, percent)
		}
	}
	_ = scanner.Err()

	if waitError := command.Wait(); waitError != nil {
		scrubbed := strings.ReplaceAll(captured.String(), token, RedactedToken)
		return fmt.Errorf(FetchFailed, waitError, scrubbed)
	}

	return nil
}
