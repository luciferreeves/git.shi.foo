package auth

import (
	"context"
	"strconv"

	"git.shi.foo/services/credentials"
	"git.shi.foo/utils/github"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func CompleteLogin(userContext context.Context, code string) (string, *fiber.Error) {
	token, exchangeError := github.Exchange(userContext, code)
	if exchangeError != nil {
		logger.Errorf(LogPrefix, TokenExchangeLog, exchangeError)
		return "", shortcuts.ServiceError(fiber.StatusBadGateway, TokenExchangeFailed)
	}

	identity, identityError := github.FetchIdentity(userContext, token)
	if identityError != nil {
		logger.Errorf(LogPrefix, IdentityFetchLog, identityError)
		return "", shortcuts.ServiceError(fiber.StatusBadGateway, IdentityFetchFailed)
	}

	if identity.Email == "" {
		if primaryEmail, emailError := github.FetchPrimaryEmail(userContext, token); emailError != nil {
			logger.Errorf(LogPrefix, EmailFetchLog, emailError)
		} else {
			identity.Email = primaryEmail
		}
	}

	providerID := strconv.FormatInt(identity.ID, 10)

	storedUser, admitError := admitUser(providerID, identity)
	if admitError != nil {
		return "", admitError
	}

	if storeError := storeRefreshToken(storedUser.ID, token.RefreshToken); storeError != nil {
		logger.Errorf(LogPrefix, CredentialStoreLog, storeError)
		return "", shortcuts.ServiceError(fiber.StatusInternalServerError, CredentialStoreFailed)
	}

	credentials.SeedAccessToken(storedUser.ID, token.AccessToken, token.Expiry)

	return providerID, nil
}
