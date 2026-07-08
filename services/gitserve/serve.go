package gitserve

import (
	"context"

	"git.shi.foo/git"
	"git.shi.foo/services/credentials"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func Advertise(service string, owner string, name string) ([]byte, *fiber.Error) {
	advertisement, adviseError := git.AdvertiseRefs(service, owner, name)
	if adviseError != nil {
		logger.Errorf(LogPrefix, AdvertiseLog, adviseError)
		return nil, shortcuts.ServiceError(fiber.StatusNotFound, RepoUnavailable)
	}
	return advertisement, nil
}

func TokenFor(requestContext context.Context, userID uint) (string, *fiber.Error) {
	accessToken, tokenError := credentials.AccessTokenForUser(requestContext, userID)
	if tokenError != nil {
		logger.Errorf(LogPrefix, TokenLog, tokenError)
		return "", shortcuts.ServiceError(fiber.StatusBadGateway, ServiceUnavailable)
	}
	return accessToken, nil
}
