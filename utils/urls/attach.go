package urls

import "github.com/gofiber/fiber/v2"

func Attach(application *fiber.App) {
	registry.Mutex.Lock()
	defer registry.Mutex.Unlock()

	for _, route := range registry.Routes.All() {
		bindPath(application, route)
	}
}
