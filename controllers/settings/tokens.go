package settings

import (
	servicesettings "git.shi.foo/services/settings"
	"git.shi.foo/utils/meta"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func CreateToken(context *fiber.Ctx) error {
	currentUser := meta.User(context)
	if guardError := servicesettings.EnsureUser(currentUser); guardError != nil {
		return guardError
	}

	request, parseError := meta.Body[CreateTokenRequest](context)
	if parseError != nil {
		return shortcuts.ServiceError(fiber.StatusBadRequest, InvalidForm)
	}

	plaintext, createError := servicesettings.CreateToken(currentUser.ID, request.Label)
	if createError != nil {
		return createError
	}

	data, dataError := servicesettings.GetIndexData(currentUser)
	if dataError != nil {
		return dataError
	}
	data.NewToken = plaintext

	meta.SetPageTitle(context, data.Title)
	return shortcuts.Render(context, "settings/index", data)
}

func RevokeToken(context *fiber.Ctx) error {
	currentUser := meta.User(context)
	if guardError := servicesettings.EnsureUser(currentUser); guardError != nil {
		return guardError
	}

	tokenID, parseError := context.ParamsInt("id")
	if parseError != nil {
		return shortcuts.ServiceError(fiber.StatusBadRequest, InvalidToken)
	}

	if revokeError := servicesettings.RevokeToken(currentUser.ID, uint(tokenID)); revokeError != nil {
		return revokeError
	}

	return shortcuts.Redirect(context, "settings.index")
}
