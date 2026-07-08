package repos

import (
	"git.shi.foo/services/repos"
	"git.shi.foo/utils/meta"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func Index(context *fiber.Ctx) error {
	data, dataError := repos.GetIndexData()
	if dataError != nil {
		return dataError
	}

	meta.SetPageTitle(context, data.Title)
	return shortcuts.Render(context, "repos/index", data)
}

func ImportIndex(context *fiber.Ctx) error {
	data, dataError := repos.GetImportData(context.UserContext(), meta.User(context))
	if dataError != nil {
		return dataError
	}

	meta.SetPageTitle(context, data.Title)
	return shortcuts.Render(context, "repos/import", data)
}

func Show(context *fiber.Ctx) error {
	data, dataError := repos.GetShowData(context.UserContext(), meta.User(context), context.Params("owner"), context.Params("repo"), context.Params("*"))
	if dataError != nil {
		return dataError
	}

	meta.SetPageTitle(context, data.Title)
	return shortcuts.Render(context, "repos/show", data)
}

func Blob(context *fiber.Ctx) error {
	data, dataError := repos.GetBlobData(context.UserContext(), meta.User(context), context.Params("owner"), context.Params("repo"), context.Params("*"))
	if dataError != nil {
		return dataError
	}

	meta.SetPageTitle(context, data.Title)
	return shortcuts.Render(context, "repos/blob", data)
}
