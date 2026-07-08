package git

import (
	"strconv"
	"strings"
)

type TreeEntry struct {
	Type string
	Name string
	Size int64
}

func Tree(owner string, name string, ref string, path string) ([]TreeEntry, error) {
	treeish := ref
	if path != "" {
		treeish = ref + ":" + path
	}

	output, runError := runGit(RepoPath(owner, name), "ls-tree", "-l", treeish)
	if runError != nil {
		return nil, runError
	}

	entries := make([]TreeEntry, 0)
	for _, line := range strings.Split(strings.TrimRight(string(output), "\n"), "\n") {
		if line == "" {
			continue
		}

		tab := strings.IndexByte(line, '\t')
		if tab < 0 {
			continue
		}

		fields := strings.Fields(line[:tab])
		if len(fields) < 4 {
			continue
		}

		var size int64
		if fields[3] != "-" {
			size, _ = strconv.ParseInt(fields[3], 10, 64)
		}

		entries = append(entries, TreeEntry{
			Type: fields[1],
			Name: line[tab+1:],
			Size: size,
		})
	}

	return entries, nil
}
