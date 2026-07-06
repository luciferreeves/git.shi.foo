package admin

import (
	"git.shi.foo/account"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func EnsureAdmin(currentUser *account.Response) *fiber.Error {
	if currentUser == nil {
		return shortcuts.ServiceError(fiber.StatusUnauthorized, AuthRequired)
	}

	if !currentUser.Admin {
		return shortcuts.ServiceError(fiber.StatusForbidden, AdminOnly)
	}

	return nil
}
