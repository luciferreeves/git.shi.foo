package repos

import (
	"bytes"
	"context"
	"fmt"
	"html"

	"git.shi.foo/account"
	"git.shi.foo/git"
	"git.shi.foo/repositories/repo"
	"git.shi.foo/utils/highlight"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func GetBlobData(requestContext context.Context, currentUser *account.Response, owner string, name string, remainder string) (*BlobContext, *fiber.Error) {
	record, viewError := resolveViewable(requestContext, currentUser, owner, name)
	if viewError != nil {
		return nil, viewError
	}

	if record.Status != repo.StatusActive {
		return nil, shortcuts.ServiceError(fiber.StatusNotFound, FileNotFound)
	}

	ref, path := resolveRef(record.Owner, record.Name, remainder, record.DefaultBranch)
	if path == "" {
		return nil, shortcuts.ServiceError(fiber.StatusNotFound, FileNotFound)
	}

	content, blobError := git.Blob(record.Owner, record.Name, ref, path)
	if blobError != nil {
		return nil, shortcuts.ServiceError(fiber.StatusNotFound, FileNotFound)
	}

	filename := fileName(path)
	branchOptions, tagOptions := buildRefOptions(record.Owner, record.Name, ref, func(option string) string {
		return blobURL(record.Owner, record.Name, option, path)
	})

	blobContext := &BlobContext{
		Title:         record.Owner + "/" + record.Name,
		Owner:         record.Owner,
		Name:          record.Name,
		Private:       record.Private,
		Description:   record.Description,
		DefaultBranch: record.DefaultBranch,
		Status:        record.Status,
		CloneURL:      cloneURL(record.Owner, record.Name),
		SSHCloneURL:   sshCloneURL(record.Owner, record.Name),
		ActiveTab:     TabFiles,
		Ref:           ref,
		RefBranches:   branchOptions,
		RefTags:       tagOptions,
		RootURL:       treeURL(record.Owner, record.Name, ref, ""),
		Path:          path,
		Filename:      filename,
		Crumbs:        buildCrumbs(record.Owner, record.Name, ref, parentPath(path)),
		SizeLabel:     humanSize(int64(len(content))),
	}

	if isBinary(content) {
		blobContext.IsBinary = true
		return blobContext, nil
	}

	highlighted, highlightError := highlight.Highlight(filename, string(content))
	if highlightError != nil {
		logger.Errorf(LogPrefix, HighlightLog, highlightError)
		blobContext.Content = fmt.Sprintf(PlainFormat, html.EscapeString(string(content)))
		return blobContext, nil
	}

	blobContext.Content = highlighted
	return blobContext, nil
}

func isBinary(content []byte) bool {
	limit := min(len(content), binarySniffLimit)
	return bytes.IndexByte(content[:limit], 0) >= 0
}
