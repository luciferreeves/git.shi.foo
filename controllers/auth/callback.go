package auth

import (
	"git.shi.foo/services/auth"
	"git.shi.foo/sessions"
	"git.shi.foo/utils/github"
	"git.shi.foo/utils/meta"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func Callback(context *fiber.Ctx) error {
	session := meta.Session(context)
	if session == nil {
		return shortcuts.ServiceError(fiber.StatusBadRequest, SessionMissing)
	}

	expectedState, ok := sessions.Get(session, github.StateKey).(string)
	providedState := context.Query(github.QueryParamState)
	if !ok || expectedState == "" || expectedState != providedState {
		return shortcuts.ServiceError(fiber.StatusBadRequest, StateMismatch)
	}

	code := context.Query(github.QueryParamCode)
	if code == "" {
		return shortcuts.ServiceError(fiber.StatusBadRequest, CodeMissing)
	}

	providerID, loginError := auth.CompleteLogin(context.UserContext(), code)
	if loginError != nil {
		return loginError
	}

	if sessionError := sessions.CreateSession(session, providerID); sessionError != nil {
		return shortcuts.ServiceError(fiber.StatusInternalServerError, SessionStartFailed)
	}

	_ = sessions.Delete(session, github.StateKey)

	return shortcuts.RedirectToPath(context, "/")
}
