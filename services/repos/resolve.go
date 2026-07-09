package repos

import (
	"context"
	"errors"
	"fmt"

	"git.shi.foo/account"
	"git.shi.foo/config"
	"git.shi.foo/models"
	"git.shi.foo/repositories/repo"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func resolveViewable(requestContext context.Context, currentUser *account.Response, owner string, name string) (*models.Repo, *fiber.Error) {
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

	return record, nil
}

func cloneURL(owner string, name string) string {
	return fmt.Sprintf("%s/%s/%s.git", config.Server.PublicURL, owner, name)
}
