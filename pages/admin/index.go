package admin

import (
	"git.shi.foo/services/admin"
	"git.shi.foo/utils/meta"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func Index(context *fiber.Ctx) error {
	indexData, dataError := admin.GetIndexData(meta.User(context))
	if dataError != nil {
		return dataError
	}

	meta.SetPageTitle(context, indexData.Title)
	return shortcuts.Render(context, "admin/index", indexData)
}
