package router

import (
	"git.shi.foo/utils/shortcuts"
	"git.shi.foo/utils/urls"

	"github.com/gofiber/fiber/v2"
)

func Initialize(application *fiber.App) {
	application.Static("/static", "./static")
	urls.Attach(application)
}

func ErrorHandler(context *fiber.Ctx, err error) error {
	fiberErr, ok := err.(*fiber.Error)
	if !ok {
		fiberErr = fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return shortcuts.RouteError(context, fiberErr)
}
