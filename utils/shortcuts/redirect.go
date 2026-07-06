package shortcuts

import (
	"git.shi.foo/utils/collections"
	"git.shi.foo/utils/urls"

	"github.com/gofiber/fiber/v2"
)

func Redirect(context *fiber.Ctx, routeName string) error {
	fullPath, exists := urls.GetFullPath(routeName)
	if !exists {
		return fiber.ErrNotFound
	}

	if isHtmxRequest(context) {
		return forward(context, fullPath)
	}

	return context.Redirect(fullPath)
}

func RedirectToPath(context *fiber.Ctx, path string) error {
	if isHtmxRequest(context) {
		return forward(context, path)
	}

	return context.Redirect(path)
}

func RedirectWithStatus(context *fiber.Ctx, routeName string, statusCode int) error {
	fullPath, exists := urls.GetFullPath(routeName)
	if !exists {
		return fiber.ErrNotFound
	}

	if isHtmxRequest(context) {
		return forward(context, fullPath)
	}

	return context.Redirect(fullPath, statusCode)
}

func forward(context *fiber.Ctx, path string) error {
	context.Request().URI().SetPath(path)
	context.Request().Header.SetMethod(fiber.MethodGet)
	context.App().Handler()(context.Context())
	return nil
}

func RedirectFull(context *fiber.Ctx, routeName string) error {
	fullPath, exists := urls.GetFullPath(routeName)
	if !exists {
		return fiber.ErrNotFound
	}

	if isHtmxRequest(context) {
		context.Set("HX-Redirect", fullPath)
		return context.SendStatus(fiber.StatusNoContent)
	}

	return context.Redirect(fullPath)
}

func RedirectRoute(context *fiber.Ctx, routeName string, params collections.Record[string, string]) error {
	fullPath, exists := urls.ResolvePath(routeName, params)
	if !exists {
		return fiber.ErrNotFound
	}

	if isHtmxRequest(context) {
		return forward(context, fullPath)
	}

	return context.Redirect(fullPath)
}

func RedirectExternal(context *fiber.Ctx, url string) error {
	return context.Redirect(url)
}
