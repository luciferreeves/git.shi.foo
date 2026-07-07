package repos

import (
	"git.shi.foo/account"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func EnsureUser(currentUser *account.Response) *fiber.Error {
	if currentUser == nil || !currentUser.Enabled {
		return shortcuts.ServiceError(fiber.StatusUnauthorized, AuthRequired)
	}

	return nil
}
