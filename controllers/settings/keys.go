package settings

import (
	servicesettings "git.shi.foo/services/settings"
	"git.shi.foo/utils/meta"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func AddKey(context *fiber.Ctx) error {
	currentUser := meta.User(context)
	if guardError := servicesettings.EnsureUser(currentUser); guardError != nil {
		return guardError
	}

	request, parseError := meta.Body[AddKeyRequest](context)
	if parseError != nil {
		return shortcuts.ServiceError(fiber.StatusBadRequest, InvalidForm)
	}

	if addError := servicesettings.AddKey(currentUser.ID, request.Title, request.Key); addError != nil {
		return addError
	}

	return shortcuts.Redirect(context, "settings.index")
}

func ImportKeys(context *fiber.Ctx) error {
	currentUser := meta.User(context)
	if guardError := servicesettings.EnsureUser(currentUser); guardError != nil {
		return guardError
	}

	if importError := servicesettings.ImportKeys(context.UserContext(), currentUser.ID); importError != nil {
		return importError
	}

	return shortcuts.Redirect(context, "settings.index")
}

func RemoveKey(context *fiber.Ctx) error {
	currentUser := meta.User(context)
	if guardError := servicesettings.EnsureUser(currentUser); guardError != nil {
		return guardError
	}

	keyID, parseError := context.ParamsInt("id")
	if parseError != nil {
		return shortcuts.ServiceError(fiber.StatusBadRequest, InvalidKeyID)
	}

	if removeError := servicesettings.RemoveKey(currentUser.ID, uint(keyID)); removeError != nil {
		return removeError
	}

	return shortcuts.Redirect(context, "settings.index")
}
