package admin

import (
	adminservice "git.shi.foo/services/admin"
	"git.shi.foo/services/invites"
	"git.shi.foo/utils/meta"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func CreateInvite(context *fiber.Ctx) error {
	currentUser := meta.User(context)
	if guardError := adminservice.EnsureAdmin(currentUser); guardError != nil {
		return guardError
	}

	request, parseError := meta.Body[CreateInviteRequest](context)
	if parseError != nil {
		return shortcuts.ServiceError(fiber.StatusBadRequest, InvalidForm)
	}

	if createError := invites.CreateInvite(currentUser.ID, request.Email, request.Username); createError != nil {
		return createError
	}

	return shortcuts.Redirect(context, "admin.index")
}
