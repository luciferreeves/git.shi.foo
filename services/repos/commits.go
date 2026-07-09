package repos

import (
	"context"
	"errors"
	"strconv"

	"git.shi.foo/account"
	"git.shi.foo/git"
	"git.shi.foo/repositories/repo"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetCommitsData(requestContext context.Context, currentUser *account.Response, owner string, name string, remainder string, page int) (*CommitsContext, *fiber.Error) {
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

	if page < 1 {
		page = 1
	}

	ref, _ := resolveRef(record.Owner, record.Name, remainder, record.DefaultBranch)

	commitsContext := &CommitsContext{
		Title:         record.Owner + "/" + record.Name,
		Owner:         record.Owner,
		Name:          record.Name,
		Private:       record.Private,
		Description:   record.Description,
		DefaultBranch: record.DefaultBranch,
		Status:        record.Status,
		CloneURL:      cloneURL(record.Owner, record.Name),
		SSHCloneURL:   sshCloneURL(record.Owner, record.Name),
		ActiveTab:     TabCommits,
		Ref:           ref,
		Page:          page,
	}

	if record.Status != repo.StatusActive {
		return commitsContext, nil
	}

	commitsContext.RefBranches, commitsContext.RefTags = buildRefOptions(record.Owner, record.Name, ref, func(option string) string {
		return commitsURL(record.Owner, record.Name, option, 1)
	})

	skip := (page - 1) * CommitsPageSize
	commits, logError := git.Log(record.Owner, record.Name, ref, CommitsPageSize+1, skip)
	if logError != nil {
		return commitsContext, nil
	}

	commitsContext.HasNext = len(commits) > CommitsPageSize
	if commitsContext.HasNext {
		commits = commits[:CommitsPageSize]
	}
	commitsContext.HasPrev = page > 1

	for index := range commits {
		commitsContext.Commits = append(commitsContext.Commits, CommitView{
			Short:   commits[index].Short,
			Message: commits[index].Message,
			Author:  commits[index].Author,
			When:    commits[index].When,
		})
	}

	if commitsContext.HasPrev {
		commitsContext.PrevURL = commitsURL(record.Owner, record.Name, ref, page-1)
	}
	if commitsContext.HasNext {
		commitsContext.NextURL = commitsURL(record.Owner, record.Name, ref, page+1)
	}

	return commitsContext, nil
}

func commitsURL(owner string, name string, ref string, page int) string {
	base := "/" + owner + "/" + name + "/commits/" + ref
	if page <= 1 {
		return base
	}
	return base + "?page=" + strconv.Itoa(page)
}
