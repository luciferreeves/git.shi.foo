package shortcuts

import "github.com/gofiber/fiber/v2"

func isHtmxRequest(context *fiber.Ctx) bool {
	return context.Get("HX-Request") == "true" && context.Get("HX-Boosted") != "true"
}
