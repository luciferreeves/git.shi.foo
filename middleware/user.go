package middleware

import (
	"git.shi.foo/repositories/user"
	"git.shi.foo/sessions"
	"git.shi.foo/utils/meta"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func resolveUser(context *fiber.Ctx) error {
	sess, ok := context.Locals(sessions.SessionLocalKey).(*session.Session)
	if !ok {
		return context.Next()
	}

	providerID := sessions.GetSessionProviderID(sess)
	if providerID == "" {
		return context.Next()
	}

	foundUser, findError := user.FindByProviderID(providerID)
	if findError != nil {
		return context.Next()
	}

	context.Locals(meta.UserKey, foundUser.ToResponse())
	return context.Next()
}
