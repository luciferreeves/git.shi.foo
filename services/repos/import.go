package repos

import (
	"context"
	"time"

	"git.shi.foo/account"
	"git.shi.foo/jobs"
	"git.shi.foo/models"
	jobrepo "git.shi.foo/repositories/job"
	"git.shi.foo/repositories/repo"
	"git.shi.foo/services/credentials"
	"git.shi.foo/utils/github"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func ImportRepo(currentUser *account.Response, owner string, name string) *fiber.Error {
	if guardError := EnsureUser(currentUser); guardError != nil {
		return guardError
	}

	if _, findError := repo.FindByOwnerName(owner, name); findError == nil {
		return shortcuts.ServiceError(fiber.StatusConflict, AlreadyImported)
	}

	accessToken, tokenError := credentials.AccessTokenForUser(context.Background(), currentUser.ID)
	if tokenError != nil {
		logger.Errorf(LogPrefix, TokenLog, tokenError)
		return shortcuts.ServiceError(fiber.StatusBadGateway, ImportFailed)
	}

	metadata, metadataError := github.FetchRepo(context.Background(), accessToken, owner, name)
	if metadataError != nil {
		logger.Errorf(LogPrefix, MetadataLog, metadataError)
		return shortcuts.ServiceError(fiber.StatusBadGateway, ImportFailed)
	}

	record := &models.Repo{
		GitHubID:      metadata.ID,
		Owner:         metadata.Owner.Login,
		Name:          metadata.Name,
		Description:   metadata.Description,
		Private:       metadata.Private,
		DefaultBranch: metadata.DefaultBranch,
		Status:        repo.StatusImporting,
		ImportedBy:    currentUser.ID,
	}
	if createError := repo.Create(record); createError != nil {
		logger.Errorf(LogPrefix, CreateLog, createError)
		return shortcuts.ServiceError(fiber.StatusInternalServerError, ImportFailed)
	}

	if enqueueError := enqueueImport(record.ID); enqueueError != nil {
		logger.Errorf(LogPrefix, EnqueueLog, enqueueError)
		return shortcuts.ServiceError(fiber.StatusInternalServerError, ImportFailed)
	}

	return nil
}

func enqueueImport(repoID uint) error {
	jobRecord := &models.Job{
		Kind:        jobrepo.KindImport,
		RepoID:      &repoID,
		Status:      jobrepo.StatusPending,
		MaxAttempts: ImportMaxAttempts,
		RunAfter:    time.Now(),
	}
	if createError := jobrepo.Create(jobRecord); createError != nil {
		return createError
	}

	jobs.Signal()
	return nil
}
