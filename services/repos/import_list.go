package repos

import (
	"context"

	"git.shi.foo/account"
	"git.shi.foo/repositories/repo"
	"git.shi.foo/services/credentials"
	"git.shi.foo/utils/github"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func GetImportData(requestContext context.Context, currentUser *account.Response) (*ImportContext, *fiber.Error) {
	if guardError := EnsureUser(currentUser); guardError != nil {
		return nil, guardError
	}

	fetched, fetchError := fetchRepos(requestContext, currentUser.ID)
	if fetchError != nil {
		return nil, fetchError
	}

	mirrored, mirroredError := mirroredSet()
	if mirroredError != nil {
		logger.Errorf(LogPrefix, ListLog, mirroredError)
		return nil, shortcuts.ServiceError(fiber.StatusInternalServerError, ListFailed)
	}

	views := make([]GitHubRepoView, 0, len(fetched))
	for _, item := range fetched {
		views = append(views, GitHubRepoView{
			Owner:       item.Owner.Login,
			Name:        item.Name,
			FullName:    item.FullName,
			Private:     item.Private,
			Description: item.Description,
			Mirrored:    mirrored[item.Owner.Login+"/"+item.Name],
		})
	}

	return &ImportContext{Title: ImportTitle, Repos: views}, nil
}

func fetchRepos(requestContext context.Context, userID uint) ([]github.Repository, *fiber.Error) {
	accessToken, tokenError := credentials.AccessTokenForUser(requestContext, userID)
	if tokenError != nil {
		logger.Errorf(LogPrefix, TokenLog, tokenError)
		return nil, shortcuts.ServiceError(fiber.StatusBadGateway, ReposFetchFailed)
	}

	fetched, fetchError := github.FetchUserRepos(requestContext, accessToken)
	if fetchError == nil {
		return fetched, nil
	}

	freshToken, refreshError := credentials.RefreshAccessTokenForUser(requestContext, userID)
	if refreshError != nil {
		logger.Errorf(LogPrefix, ReposFetchLog, refreshError)
		return nil, shortcuts.ServiceError(fiber.StatusBadGateway, ReposFetchFailed)
	}

	fetched, fetchError = github.FetchUserRepos(requestContext, freshToken)
	if fetchError != nil {
		logger.Errorf(LogPrefix, ReposFetchLog, fetchError)
		return nil, shortcuts.ServiceError(fiber.StatusBadGateway, ReposFetchFailed)
	}

	return fetched, nil
}

func mirroredSet() (map[string]bool, error) {
	records, listError := repo.ListAll()
	if listError != nil {
		return nil, listError
	}

	set := make(map[string]bool, len(records))
	for _, record := range records {
		set[record.Owner+"/"+record.Name] = true
	}

	return set, nil
}
