package git

import (
	"errors"
	"strings"
	"time"
)

type Commit struct {
	SHA     string
	Short   string
	Message string
	Author  string
	When    time.Time
}

func LatestCommit(owner string, name string, ref string) (*Commit, error) {
	output, runError := runGit(RepoPath(owner, name), "log", "-1", CommitFormat, ref)
	if runError != nil {
		return nil, runError
	}
	return parseCommit(output)
}

func LastCommitForPath(owner string, name string, ref string, path string) (*Commit, error) {
	output, runError := runGit(RepoPath(owner, name), "log", "-1", CommitFormat, ref, "--", path)
	if runError != nil {
		return nil, runError
	}
	return parseCommit(output)
}

func parseCommit(output []byte) (*Commit, error) {
	trimmed := strings.TrimRight(string(output), "\n")
	if trimmed == "" {
		return nil, errors.New(CommitParseFailed)
	}

	fields := strings.Split(trimmed, "\x1f")
	if len(fields) < 5 {
		return nil, errors.New(CommitParseFailed)
	}

	when, _ := time.Parse(time.RFC3339, fields[4])

	return &Commit{
		SHA:     fields[0],
		Short:   fields[1],
		Message: fields[2],
		Author:  fields[3],
		When:    when,
	}, nil
}
