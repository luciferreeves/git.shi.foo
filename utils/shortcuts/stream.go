package shortcuts

import (
	"bufio"

	"github.com/gofiber/fiber/v2"
)

func SSE(context *fiber.Ctx, handler func(*bufio.Writer)) error {
	context.Set("Content-Type", "text/event-stream")
	context.Set("Cache-Control", "no-cache")
	context.Set("Connection", "keep-alive")
	context.Set("X-Accel-Buffering", "no")

	context.Context().SetBodyStreamWriter(handler)

	return nil
}
