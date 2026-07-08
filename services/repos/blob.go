package repos

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html"

	"git.shi.foo/account"
	"git.shi.foo/git"
	"git.shi.foo/repositories/repo"
	"git.shi.foo/utils/highlight"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetBlobData(requestContext context.Context, currentUser *account.Response, owner string, name string, path string) (*BlobContext, *fiber.Error) {
	record, findError := repo.FindByOwnerName(owner, name)
	if findError != nil {
		if errors.Is(findError, gorm.ErrRecordNotFound) {
			return nil, shortcuts.ServiceError(fiber.StatusNotFound, RepoNotFound)
		}
		logger.Errorf(LogPrefix, ListLog, findError)
		return nil, shortcuts.ServiceError(fiber.StatusInternalServerError, ListFailed)
	}

	if viewError := ensureViewable(requestContext, currentUser, record); viewError != nil {
		return nil, viewError
	}

	if record.Status != repo.StatusActive {
		return nil, shortcuts.ServiceError(fiber.StatusNotFound, FileNotFound)
	}

	content, blobError := git.Blob(record.Owner, record.Name, git.HeadRef, path)
	if blobError != nil {
		return nil, shortcuts.ServiceError(fiber.StatusNotFound, FileNotFound)
	}

	filename := fileName(path)
	blobContext := &BlobContext{
		Title:     record.Owner + "/" + record.Name,
		Owner:     record.Owner,
		Name:      record.Name,
		Path:      path,
		Filename:  filename,
		Crumbs:    buildCrumbs(record.Owner, record.Name, parentPath(path)),
		SizeLabel: humanSize(int64(len(content))),
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
