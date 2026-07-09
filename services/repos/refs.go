package repos

import (
	"strings"

	"git.shi.foo/git"
)

func resolveRef(owner string, name string, remainder string, defaultBranch string) (string, string) {
	fallback := defaultBranch
	if fallback == "" {
		fallback = git.HeadRef
	}

	if remainder == "" {
		return fallback, ""
	}

	segments := strings.Split(remainder, "/")

	if names, namesError := git.RefNames(owner, name); namesError == nil {
		for count := len(segments); count > 0; count-- {
			candidate := strings.Join(segments[:count], "/")
			if _, present := names[candidate]; present {
				return candidate, strings.Join(segments[count:], "/")
			}
		}
	}

	if isNamedRef(segments[0]) && git.VerifyRef(owner, name, segments[0]) {
		return segments[0], strings.Join(segments[1:], "/")
	}

	return fallback, remainder
}

func isNamedRef(candidate string) bool {
	return candidate != "" && !strings.HasPrefix(candidate, "-")
}

func buildRefOptions(owner string, name string, current string, path string, blob bool) ([]RefOption, []RefOption) {
	collect := func(refs []git.Ref, isTag bool) []RefOption {
		options := make([]RefOption, 0, len(refs))
		for _, ref := range refs {
			target := treeURL(owner, name, ref.Name, path)
			if blob {
				target = blobURL(owner, name, ref.Name, path)
			}
			options = append(options, RefOption{
				Name:   ref.Name,
				URL:    target,
				Active: ref.Name == current,
				IsTag:  isTag,
			})
		}
		return options
	}

	branchOptions := make([]RefOption, 0)
	tagOptions := make([]RefOption, 0)

	if branches, branchError := git.Branches(owner, name); branchError == nil {
		branchOptions = collect(branches, false)
	}
	if tags, tagError := git.Tags(owner, name); tagError == nil {
		tagOptions = collect(tags, true)
	}

	return branchOptions, tagOptions
}
