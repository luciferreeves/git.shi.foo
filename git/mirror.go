package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var progressPattern = regexp.MustCompile(`(Counting objects|Compressing objects|Receiving objects|Resolving deltas):\s+(\d+)%`)

func MirrorWithProgress(owner string, name string, token string, onProgress func(phase string, percent int)) error {
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
		"clone", "--mirror", "--progress", cleanURL, destination,
	)
	cloneCommand.Env = append(os.Environ(), GitTokenEnv+"="+token)

	stderrPipe, pipeError := cloneCommand.StderrPipe()
	if pipeError != nil {
		return fmt.Errorf(CloneFailed, pipeError, "")
	}

	if startError := cloneCommand.Start(); startError != nil {
		return fmt.Errorf(CloneFailed, startError, "")
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

	if waitError := cloneCommand.Wait(); waitError != nil {
		scrubbed := strings.ReplaceAll(captured.String(), token, RedactedToken)
		return fmt.Errorf(CloneFailed, waitError, scrubbed)
	}

	return nil
}

func scanProgressLines(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	for index := range len(data) {
		if data[index] == '\n' || data[index] == '\r' {
			return index + 1, data[:index], nil
		}
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

func parseProgress(line string) (string, int, bool) {
	match := progressPattern.FindStringSubmatch(line)
	if match == nil {
		return "", 0, false
	}

	percent, convertError := strconv.Atoi(match[2])
	if convertError != nil {
		return "", 0, false
	}

	return progressPhase(match[1]), percent, true
}

func progressPhase(label string) string {
	switch label {
	case "Counting objects":
		return PhaseCounting
	case "Compressing objects":
		return PhaseCompressing
	case "Receiving objects":
		return PhaseReceiving
	case "Resolving deltas":
		return PhaseResolving
	default:
		return label
	}
}
