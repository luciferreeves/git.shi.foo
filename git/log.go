package git

import (
	"strconv"
	"strings"
	"time"
)

func Log(owner string, name string, ref string, limit int, skip int) ([]Commit, error) {
	output, runError := runGit(RepoPath(owner, name), "log", CommitFormat, "-n", strconv.Itoa(limit), "--skip", strconv.Itoa(skip), ref)
	if runError != nil {
		return nil, runError
	}

	commits := make([]Commit, 0)
	for line := range strings.SplitSeq(strings.TrimRight(string(output), "\n"), "\n") {
		if line == "" {
			continue
		}

		fields := strings.Split(line, "\x1f")
		if len(fields) < 5 {
			continue
		}

		when, _ := time.Parse(time.RFC3339, fields[4])
		commits = append(commits, Commit{
			SHA:     fields[0],
			Short:   fields[1],
			Message: fields[2],
			Author:  fields[3],
			When:    when,
		})
	}

	return commits, nil
}
