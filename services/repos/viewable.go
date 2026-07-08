package repos

import (
	"context"
	"errors"

	"git.shi.foo/account"
	"git.shi.foo/repositories/repo"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Viewable(requestContext context.Context, currentUser *account.Response, owner string, name string) *fiber.Error {
	record, findError := repo.FindByOwnerName(owner, name)
	if findError != nil {
		if errors.Is(findError, gorm.ErrRecordNotFound) {
			return shortcuts.ServiceError(fiber.StatusNotFound, RepoNotFound)
		}
		logger.Errorf(LogPrefix, ListLog, findError)
		return shortcuts.ServiceError(fiber.StatusInternalServerError, ListFailed)
	}

	return ensureViewable(requestContext, currentUser, record)
}
