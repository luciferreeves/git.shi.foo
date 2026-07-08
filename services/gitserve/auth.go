package gitserve

import (
	"context"

	"git.shi.foo/account"
	"git.shi.foo/git"
	"git.shi.foo/repositories/repo"
	patrepo "git.shi.foo/repositories/token"
	"git.shi.foo/repositories/user"
	"git.shi.foo/services/credentials"
	"git.shi.foo/services/repos"
	"git.shi.foo/utils/github"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"
	"git.shi.foo/utils/tokens"

	"github.com/gofiber/fiber/v2"
)

func ResolveUser(secret string) *account.Response {
	if secret == "" {
		return nil
	}

	patRecord, findError := patrepo.FindByHash(tokens.Hash(secret))
	if findError != nil {
		return nil
	}

	userRecord, userError := user.FindByID(patRecord.UserID)
	if userError != nil {
		return nil
	}

	response := userRecord.ToResponse()
	return &response
}

func Authorize(requestContext context.Context, currentUser *account.Response, owner string, name string, service string) *fiber.Error {
	if service == git.ServiceReceivePack {
		if currentUser == nil {
			return shortcuts.ServiceError(fiber.StatusUnauthorized, AuthPrompt)
		}
		if !currentUser.Enabled {
			return shortcuts.ServiceError(fiber.StatusForbidden, NotAllowed)
		}
		return authorizePush(requestContext, currentUser, owner, name)
	}

	record, findError := repo.FindByOwnerName(owner, name)
	if findError != nil {
		return shortcuts.ServiceError(fiber.StatusNotFound, RepoNotFound)
	}
	if !record.Private {
		return nil
	}
	if currentUser == nil {
		return shortcuts.ServiceError(fiber.StatusUnauthorized, AuthPrompt)
	}

	return repos.Viewable(requestContext, currentUser, owner, name)
}

func authorizePush(requestContext context.Context, currentUser *account.Response, owner string, name string) *fiber.Error {
	accessToken, tokenError := credentials.AccessTokenForUser(requestContext, currentUser.ID)
	if tokenError != nil {
		logger.Errorf(LogPrefix, TokenLog, tokenError)
		return shortcuts.ServiceError(fiber.StatusBadGateway, ServiceUnavailable)
	}

	metadata, fetchError := github.FetchRepo(requestContext, accessToken, owner, name)
	if fetchError != nil {
		return shortcuts.ServiceError(fiber.StatusNotFound, RepoNotFound)
	}

	if !metadata.Permissions.Push {
		return shortcuts.ServiceError(fiber.StatusForbidden, NotAllowed)
	}

	return nil
}
