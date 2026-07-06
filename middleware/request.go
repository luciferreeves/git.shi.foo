package middleware

import (
	"git.shi.foo/utils/meta"

	"github.com/gofiber/fiber/v2"
)

func request(context *fiber.Ctx) error {
	context.Locals(meta.RequestKey, meta.BuildRequest(context))
	return context.Next()
}
