package repos

import (
	servicerepos "git.shi.foo/services/repos"
	"git.shi.foo/utils/meta"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func Retry(context *fiber.Ctx) error {
	owner := context.Params("owner")
	name := context.Params("repo")

	if retryError := servicerepos.RetryImport(context.UserContext(), meta.User(context), owner, name); retryError != nil {
		return retryError
	}

	return shortcuts.RedirectToPath(context, "/"+owner+"/"+name)
}
