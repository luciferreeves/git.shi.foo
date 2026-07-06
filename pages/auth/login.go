package auth

import (
	"git.shi.foo/services/auth"
	"git.shi.foo/sessions"
	"git.shi.foo/utils/github"
	"git.shi.foo/utils/meta"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func Login(context *fiber.Ctx) error {
	authURL, state, buildError := auth.BuildAuthURL()
	if buildError != nil {
		return buildError
	}

	session := meta.Session(context)
	if session == nil {
		return shortcuts.ServiceError(fiber.StatusInternalServerError, SessionMissing)
	}

	if storeError := sessions.Set(session, github.StateKey, state); storeError != nil {
		return shortcuts.ServiceError(fiber.StatusInternalServerError, StateStoreFailed)
	}

	return shortcuts.RedirectExternal(context, authURL)
}
