package git

import (
	"strings"
	"time"
)

type Ref struct {
	Name    string
	Short   string
	Message string
	Author  string
	When    time.Time
}

func Branches(owner string, name string) ([]Ref, error) {
	return refs(owner, name, RefsHeads)
}

func Tags(owner string, name string) ([]Ref, error) {
	return refs(owner, name, RefsTags)
}

func RefNames(owner string, name string) (map[string]struct{}, error) {
	names := make(map[string]struct{})

	for _, prefix := range []string{RefsHeads, RefsTags} {
		output, runError := runGit(RepoPath(owner, name), "for-each-ref", "--format=%(refname:short)", prefix)
		if runError != nil {
			return nil, runError
		}

		for line := range strings.SplitSeq(strings.TrimRight(string(output), "\n"), "\n") {
			if line != "" {
				names[line] = struct{}{}
			}
		}
	}

	return names, nil
}

func VerifyRef(owner string, name string, ref string) bool {
	_, runError := runGit(RepoPath(owner, name), "rev-parse", "--verify", "--quiet", ref+"^{commit}")
	return runError == nil
}

func refs(owner string, name string, prefix string) ([]Ref, error) {
	output, runError := runGit(RepoPath(owner, name), "for-each-ref", "--sort=-creatordate", "--format="+RefFormat, prefix)
	if runError != nil {
		return nil, runError
	}

	found := make([]Ref, 0)
	for line := range strings.SplitSeq(strings.TrimRight(string(output), "\n"), "\n") {
		if line == "" {
			continue
		}

		fields := strings.Split(line, "\x1f")
		if len(fields) < 5 {
			continue
		}

		when, _ := time.Parse(time.RFC3339, fields[4])
		found = append(found, Ref{
			Name:    fields[0],
			Short:   fields[1],
			Message: fields[2],
			Author:  fields[3],
			When:    when,
		})
	}

	return found, nil
}
