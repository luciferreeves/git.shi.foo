package repos

import (
	"context"
	"errors"
	"fmt"

	"git.shi.foo/account"
	"git.shi.foo/config"
	"git.shi.foo/git"
	"git.shi.foo/repositories/repo"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetShowData(requestContext context.Context, currentUser *account.Response, owner string, name string, path string) (*ShowContext, *fiber.Error) {
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

	showContext := &ShowContext{
		Title:         record.Owner + "/" + record.Name,
		Owner:         record.Owner,
		Name:          record.Name,
		Private:       record.Private,
		Description:   record.Description,
		DefaultBranch: record.DefaultBranch,
		Status:        record.Status,
		CloneURL:      fmt.Sprintf("%s/%s/%s.git", config.Server.PublicURL, record.Owner, record.Name),
		Ready:         record.Status == repo.StatusActive,
		Path:          path,
		Crumbs:        buildCrumbs(record.Owner, record.Name, path),
	}

	if !showContext.Ready {
		return showContext, nil
	}

	if entries, treeError := git.Tree(record.Owner, record.Name, git.HeadRef, path); treeError == nil {
		showContext.Entries = toEntryViews(entries, record.Owner, record.Name, path)
	}

	if commit, commitError := git.LatestCommit(record.Owner, record.Name, git.HeadRef); commitError == nil {
		showContext.LatestCommit = &CommitView{
			Short:   commit.Short,
			Message: commit.Message,
			Author:  commit.Author,
			When:    commit.When,
		}
	}

	return showContext, nil
}

func toEntryViews(entries []git.TreeEntry, owner string, name string, path string) []EntryView {
	directories := make([]EntryView, 0)
	files := make([]EntryView, 0)

	for _, entry := range entries {
		full := childPath(path, entry.Name)
		view := EntryView{
			Type:      entry.Type,
			Name:      entry.Name,
			Size:      entry.Size,
			SizeLabel: DirSizeLabel,
			IsDir:     entry.Type == git.TypeTree,
		}
		if view.IsDir {
			view.URL = treeURL(owner, name, full)
			directories = append(directories, view)
		} else {
			view.URL = blobURL(owner, name, full)
			view.SizeLabel = humanSize(entry.Size)
			files = append(files, view)
		}
	}

	return append(directories, files...)
}

func humanSize(size int64) string {
	switch {
	case size >= bytesPerMB:
		return fmt.Sprintf(SizeFormatMB, float64(size)/bytesPerMB)
	case size >= bytesPerKB:
		return fmt.Sprintf(SizeFormatKB, float64(size)/bytesPerKB)
	default:
		return fmt.Sprintf(SizeFormatB, size)
	}
}
