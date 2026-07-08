package repos

import (
	"context"

	"git.shi.foo/account"
	jobrepo "git.shi.foo/repositories/job"
	"git.shi.foo/repositories/repo"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func RetryImport(requestContext context.Context, currentUser *account.Response, owner string, name string) *fiber.Error {
	if guardError := EnsureUser(currentUser); guardError != nil {
		return guardError
	}

	record, findError := repo.FindByOwnerName(owner, name)
	if findError != nil {
		return shortcuts.ServiceError(fiber.StatusNotFound, RepoNotFound)
	}

	if viewError := ensureViewable(requestContext, currentUser, record); viewError != nil {
		return viewError
	}

	if record.ImportedBy != currentUser.ID && !currentUser.Admin {
		return shortcuts.ServiceError(fiber.StatusForbidden, RetryForbidden)
	}

	open, checkError := jobrepo.HasOpenForRepo(record.ID, jobrepo.KindImport)
	if checkError != nil {
		logger.Errorf(LogPrefix, EnqueueLog, checkError)
		return shortcuts.ServiceError(fiber.StatusInternalServerError, ImportFailed)
	}
	if open {
		return nil
	}

	record.Status = repo.StatusImporting
	if updateError := repo.Update(record); updateError != nil {
		logger.Errorf(LogPrefix, StatusUpdateLog, updateError)
		return shortcuts.ServiceError(fiber.StatusInternalServerError, ImportFailed)
	}

	if enqueueError := enqueueImport(record.ID); enqueueError != nil {
		logger.Errorf(LogPrefix, EnqueueLog, enqueueError)
		return shortcuts.ServiceError(fiber.StatusInternalServerError, ImportFailed)
	}

	return nil
}
