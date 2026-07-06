package shortcuts

import (
	"git.shi.foo/utils/meta"

	"github.com/gofiber/fiber/v2"
)

func RouteError(context *fiber.Ctx, err *fiber.Error) error {
	meta.SetPageTitle(context, err.Error())
	return RenderWithStatus(context, "errors/error", err, err.Code)
}

func ServiceError(code int, message string) *fiber.Error {
	return fiber.NewError(code, message)
}
