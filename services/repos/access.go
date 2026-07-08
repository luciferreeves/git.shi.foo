package repos

import (
	"context"
	"strings"

	"git.shi.foo/account"
	"git.shi.foo/models"
	"git.shi.foo/services/credentials"
	"git.shi.foo/utils/github"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func ensureViewable(requestContext context.Context, currentUser *account.Response, record *models.Repo) *fiber.Error {
	if !record.Private {
		return nil
	}

	if currentUser == nil || !currentUser.Enabled {
		return shortcuts.ServiceError(fiber.StatusNotFound, RepoNotFound)
	}

	if record.ImportedBy == currentUser.ID {
		return nil
	}

	accessToken, tokenError := credentials.AccessTokenForUser(requestContext, currentUser.ID)
	if tokenError != nil {
		logger.Errorf(LogPrefix, TokenLog, tokenError)
		return shortcuts.ServiceError(fiber.StatusNotFound, RepoNotFound)
	}

	if _, accessError := github.FetchRepo(requestContext, accessToken, record.Owner, record.Name); accessError != nil {
		return shortcuts.ServiceError(fiber.StatusNotFound, RepoNotFound)
	}

	return nil
}

func treeURL(owner string, name string, path string) string {
	if path == "" {
		return "/" + owner + "/" + name
	}
	return "/" + owner + "/" + name + "/tree/" + path
}

func blobURL(owner string, name string, path string) string {
	return "/" + owner + "/" + name + "/blob/" + path
}

func childPath(parent string, name string) string {
	if parent == "" {
		return name
	}
	return parent + "/" + name
}

func parentPath(path string) string {
	if index := strings.LastIndexByte(path, '/'); index >= 0 {
		return path[:index]
	}
	return ""
}

func fileName(path string) string {
	if index := strings.LastIndexByte(path, '/'); index >= 0 {
		return path[index+1:]
	}
	return path
}

func buildCrumbs(owner string, name string, path string) []Crumb {
	crumbs := make([]Crumb, 0)
	if path == "" {
		return crumbs
	}

	accumulated := ""
	for _, segment := range strings.Split(path, "/") {
		accumulated = childPath(accumulated, segment)
		crumbs = append(crumbs, Crumb{Name: segment, URL: treeURL(owner, name, accumulated)})
	}

	return crumbs
}
