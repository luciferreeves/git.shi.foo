package home

import (
	"git.shi.foo/services/home"
	"git.shi.foo/utils/meta"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func Index(context *fiber.Ctx) error {
	indexData := home.GetIndexData(meta.User(context))
	meta.SetPageTitle(context, indexData.Title)
	return shortcuts.Render(context, "home/index", indexData)
}
