package middleware

import (
	"git.shi.foo/sessions"

	"github.com/gofiber/fiber/v2"
)

func userSession(context *fiber.Ctx) error {
	sess, sessionError := sessions.Store.Get(context)
	if sessionError != nil {
		return context.Next()
	}

	context.Locals(sessions.SessionLocalKey, sess)
	return context.Next()
}
