package repos

import (
	"html"
	"strings"

	"git.shi.foo/git"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/markdown"
	"git.shi.foo/utils/rst"
)

func loadReadme(owner string, name string, entries []git.TreeEntry) (string, string) {
	readmeName := pickReadme(entries)
	if readmeName == "" {
		return "", ""
	}

	content, blobError := git.Blob(owner, name, git.HeadRef, readmeName)
	if blobError != nil {
		return "", ""
	}

	switch {
	case isMarkdownName(readmeName):
		return readmeName, renderOrPlain(content, markdown.Render)
	case isRSTName(readmeName):
		return readmeName, renderOrPlain(content, rst.Render)
	default:
		return readmeName, plainReadme(content)
	}
}

func renderOrPlain(content []byte, renderer func([]byte) (string, error)) string {
	rendered, renderError := renderer(content)
	if renderError != nil {
		logger.Errorf(LogPrefix, ReadmeLog, renderError)
		return plainReadme(content)
	}
	return rendered
}

func pickReadme(entries []git.TreeEntry) string {
	fallback := ""
	for _, entry := range entries {
		if entry.Type != git.TypeBlob {
			continue
		}

		lower := strings.ToLower(entry.Name)
		if !strings.HasPrefix(lower, "readme") {
			continue
		}
		if lower == "readme.md" {
			return entry.Name
		}
		if fallback == "" {
			fallback = entry.Name
		}
	}
	return fallback
}

func isMarkdownName(name string) bool {
	lower := strings.ToLower(name)
	return strings.HasSuffix(lower, ".md") || strings.HasSuffix(lower, ".markdown")
}

func isRSTName(name string) bool {
	return strings.HasSuffix(strings.ToLower(name), ".rst")
}

func plainReadme(content []byte) string {
	return ReadmePlainOpen + html.EscapeString(string(content)) + ReadmePlainClose
}
