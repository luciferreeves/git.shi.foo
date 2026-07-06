package meta

import (
	"git.shi.foo/account"

	"github.com/gofiber/fiber/v2"
)

func User(context *fiber.Ctx) *account.Response {
	user, ok := context.Locals(UserKey).(account.Response)
	if !ok {
		return nil
	}

	return &user
}
