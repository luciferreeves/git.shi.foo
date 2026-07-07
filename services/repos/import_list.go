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

	collected, collectError := collectRepos(requestContext, accessToken)
	if collectError == nil {
		return collected, nil
	}

	freshToken, refreshError := credentials.RefreshAccessTokenForUser(requestContext, userID)
	if refreshError != nil {
		logger.Errorf(LogPrefix, ReposFetchLog, refreshError)
		return nil, shortcuts.ServiceError(fiber.StatusBadGateway, ReposFetchFailed)
	}

	collected, collectError = collectRepos(requestContext, freshToken)
	if collectError != nil {
		logger.Errorf(LogPrefix, ReposFetchLog, collectError)
		return nil, shortcuts.ServiceError(fiber.StatusBadGateway, ReposFetchFailed)
	}

	return collected, nil
}

func collectRepos(requestContext context.Context, accessToken string) ([]github.Repository, error) {
	seen := make(map[int64]bool)
	collected := make([]github.Repository, 0)

	installations, installationsError := github.FetchUserInstallations(requestContext, accessToken)
	if installationsError != nil {
		return nil, installationsError
	}

	for _, installation := range installations {
		repositories, reposError := github.FetchInstallationRepos(requestContext, accessToken, installation.ID)
		if reposError != nil {
			return nil, reposError
		}

		for _, repository := range repositories {
			if !seen[repository.ID] {
				seen[repository.ID] = true
				collected = append(collected, repository)
			}
		}
	}

	userRepos, userReposError := github.FetchUserRepos(requestContext, accessToken)
	if userReposError != nil {
		return nil, userReposError
	}

	for _, repository := range userRepos {
		if !seen[repository.ID] {
			seen[repository.ID] = true
			collected = append(collected, repository)
		}
	}

	return collected, nil
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
