package git

import (
	"sort"
	"strings"
)

type Contributor struct {
	Name  string
	Email string
	Count int
}

func Contributors(owner string, name string, ref string) ([]Contributor, error) {
	output, runError := runGit(RepoPath(owner, name), "log", ContributorFormat, ref)
	if runError != nil {
		return nil, runError
	}

	byKey := make(map[string]*Contributor)
	order := make([]string, 0)

	for line := range strings.SplitSeq(strings.TrimRight(string(output), "\n"), "\n") {
		if line == "" {
			continue
		}

		fields := strings.SplitN(line, "\x1f", 2)
		author := fields[0]
		email := ""
		if len(fields) > 1 {
			email = fields[1]
		}

		key := strings.ToLower(email)
		if key == "" {
			key = strings.ToLower(author)
		}

		if existing, present := byKey[key]; present {
			existing.Count++
			continue
		}
		byKey[key] = &Contributor{Name: author, Email: email, Count: 1}
		order = append(order, key)
	}

	contributors := make([]Contributor, 0, len(order))
	for _, key := range order {
		contributors = append(contributors, *byKey[key])
	}
	sort.SliceStable(contributors, func(first int, second int) bool {
		return contributors[first].Count > contributors[second].Count
	})

	return contributors, nil
}
