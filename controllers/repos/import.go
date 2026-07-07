package repos

import (
	servicerepos "git.shi.foo/services/repos"
	"git.shi.foo/utils/meta"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func Import(context *fiber.Ctx) error {
	request, parseError := meta.Body[ImportRequest](context)
	if parseError != nil {
		return shortcuts.ServiceError(fiber.StatusBadRequest, InvalidForm)
	}

	if importError := servicerepos.ImportRepo(meta.User(context), request.Owner, request.Name); importError != nil {
		return importError
	}

	return shortcuts.Redirect(context, "repos.index")
}
