package middleware

import (
	"git.shi.foo/config"
	"git.shi.foo/utils/meta"

	"github.com/gofiber/fiber/v2"
)

func chrome(context *fiber.Ctx) error {
	context.Locals(meta.VersionKey, config.AppVersion)
	return context.Next()
}
