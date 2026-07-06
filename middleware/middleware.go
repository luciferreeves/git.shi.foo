package middleware

import "github.com/gofiber/fiber/v2"

func Initialize(application *fiber.App) {
	application.Use(userSession)
	application.Use(request)
}
