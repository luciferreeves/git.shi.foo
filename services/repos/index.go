package repos

import (
	"git.shi.foo/models"
	"git.shi.foo/repositories/repo"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func GetIndexData() (*IndexContext, *fiber.Error) {
	records, listError := repo.ListAll()
	if listError != nil {
		logger.Errorf(LogPrefix, ListLog, listError)
		return nil, shortcuts.ServiceError(fiber.StatusInternalServerError, ListFailed)
	}

	return &IndexContext{
		Title: ReposTitle,
		Repos: toRepoViews(records),
	}, nil
}

func toRepoViews(records []models.Repo) []RepoView {
	views := make([]RepoView, 0, len(records))
	for _, record := range records {
		views = append(views, RepoView{
			Owner:       record.Owner,
			Name:        record.Name,
			Private:     record.Private,
			Description: record.Description,
			Status:      record.Status,
			UpdatedAt:   record.UpdatedAt,
		})
	}

	return views
}
