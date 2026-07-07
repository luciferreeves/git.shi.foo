package settings

import (
	"git.shi.foo/services/settings"
	"git.shi.foo/utils/meta"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func Index(context *fiber.Ctx) error {
	data, dataError := settings.GetIndexData(meta.User(context))
	if dataError != nil {
		return dataError
	}

	meta.SetPageTitle(context, data.Title)
	return shortcuts.Render(context, "settings/index", data)
}
