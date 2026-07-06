package invites

import (
	"git.shi.foo/services/invites"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func Accept(context *fiber.Ctx) error {
	token := context.Query(AcceptQueryToken)
	if acceptError := invites.AcceptInvite(token); acceptError != nil {
		return acceptError
	}

	return shortcuts.RedirectToPath(context, "/auth/login")
}
