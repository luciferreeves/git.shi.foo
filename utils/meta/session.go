package meta

import (
	"git.shi.foo/sessions"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func Session(context *fiber.Ctx) *session.Session {
	sess, ok := context.Locals(sessions.SessionLocalKey).(*session.Session)
	if !ok {
		return nil
	}

	return sess
}
