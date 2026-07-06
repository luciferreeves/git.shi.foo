package auth

import (
	"git.shi.foo/utils/github"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

func BuildAuthURL() (string, string, *fiber.Error) {
	state, stateError := github.GenerateState()
	if stateError != nil {
		logger.Errorf(LogPrefix, StateGenerationLog, stateError)
		return "", "", shortcuts.ServiceError(fiber.StatusInternalServerError, StateGenerationFailed)
	}

	return github.OAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline), state, nil
}
