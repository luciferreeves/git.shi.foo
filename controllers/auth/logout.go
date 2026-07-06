package auth

import (
	"git.shi.foo/sessions"
	"git.shi.foo/utils/meta"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func Logout(context *fiber.Ctx) error {
	session := meta.Session(context)
	if session != nil {
		_ = sessions.DestroySession(session)
	}

	return shortcuts.RedirectToPath(context, "/")
}
