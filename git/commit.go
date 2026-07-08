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
	output, runError := runGit(RepoPath(owner, name), "log", "-1", "--format=%H%x1f%h%x1f%s%x1f%an%x1f%aI", ref)
	if runError != nil {
		return nil, runError
	}

	fields := strings.Split(strings.TrimRight(string(output), "\n"), "\x1f")
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
