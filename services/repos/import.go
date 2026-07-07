package repos

import (
	"context"

	"git.shi.foo/account"
	"git.shi.foo/git"
	"git.shi.foo/models"
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

	go performMirror(record.ID, currentUser.ID, record.Owner, record.Name)

	return nil
}

func performMirror(repoID uint, userID uint, owner string, name string) {
	backgroundContext := context.Background()

	accessToken, tokenError := credentials.AccessTokenForUser(backgroundContext, userID)
	if tokenError != nil {
		logger.Errorf(LogPrefix, TokenLog, tokenError)
		setStatus(repoID, repo.StatusFailed)
		return
	}

	if mirrorError := git.Mirror(owner, name, accessToken); mirrorError != nil {
		logger.Errorf(LogPrefix, MirrorLog, mirrorError)
		setStatus(repoID, repo.StatusFailed)
		return
	}

	setStatus(repoID, repo.StatusActive)
}

func setStatus(repoID uint, status string) {
	record, findError := repo.FindByID(repoID)
	if findError != nil {
		return
	}

	record.Status = status
	_ = repo.Update(record)
}
