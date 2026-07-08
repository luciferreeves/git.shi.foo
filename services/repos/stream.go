package repos

import (
	"context"
	"errors"

	"git.shi.foo/account"
	"git.shi.foo/jobs"
	jobrepo "git.shi.foo/repositories/job"
	"git.shi.foo/repositories/repo"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ImportStream(requestContext context.Context, currentUser *account.Response, owner string, name string) (*ImportStreamContext, *fiber.Error) {
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

	streamContext := &ImportStreamContext{RepoID: record.ID}
	if latest, jobError := jobrepo.FindLatestForRepo(record.ID, jobrepo.KindImport); jobError == nil {
		streamContext.Topic = jobs.Topic(latest.ID)
	}

	return streamContext, nil
}

func ImportSnapshot(repoID uint) jobs.Progress {
	latest, jobError := jobrepo.FindLatestForRepo(repoID, jobrepo.KindImport)
	if jobError == nil {
		return jobs.Progress{
			JobID:   latest.ID,
			Status:  latest.Status,
			Phase:   latest.Phase,
			Percent: latest.Percent,
			Done:    jobs.IsTerminal(latest.Status),
		}
	}

	record, findError := repo.FindByID(repoID)
	if findError != nil {
		return jobs.Progress{Status: jobrepo.StatusFailed, Done: true}
	}

	return jobs.Progress{
		Status:  repoStatusToJob(record.Status),
		Percent: repoDonePercent(record.Status),
		Done:    record.Status != repo.StatusImporting,
	}
}

func repoStatusToJob(status string) string {
	switch status {
	case repo.StatusActive:
		return jobrepo.StatusSucceeded
	case repo.StatusFailed:
		return jobrepo.StatusFailed
	default:
		return jobrepo.StatusRunning
	}
}

func repoDonePercent(status string) int {
	if status == repo.StatusActive {
		return DonePercent
	}
	return 0
}
