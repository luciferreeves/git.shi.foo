package meta

import "github.com/gofiber/fiber/v2"

func SetPageTitle(context *fiber.Ctx, title string) {
	context.Locals("Title", title)
}
