package sync

import (
	servicesync "git.shi.foo/services/sync"

	"github.com/gofiber/fiber/v2"
)

func Webhook(context *fiber.Ctx) error {
	eventType := context.Get(GitHubEventHeader)
	signature := context.Get(SignatureHeader)

	if ingestError := servicesync.Ingest(eventType, signature, context.Body()); ingestError != nil {
		return ingestError
	}

	return context.SendStatus(fiber.StatusOK)
}
