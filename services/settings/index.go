package settings

import (
	"git.shi.foo/account"
	"git.shi.foo/models"
	"git.shi.foo/repositories/token"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func GetIndexData(currentUser *account.Response) (*IndexContext, *fiber.Error) {
	if guardError := EnsureUser(currentUser); guardError != nil {
		return nil, guardError
	}

	records, listError := token.ListByUser(currentUser.ID)
	if listError != nil {
		logger.Errorf(LogPrefix, TokenListLog, listError)
		return nil, shortcuts.ServiceError(fiber.StatusInternalServerError, TokenListFailed)
	}

	return &IndexContext{
		Title:  SettingsTitle,
		Tokens: toTokenViews(records),
	}, nil
}

func toTokenViews(records []models.PersonalAccessToken) []TokenView {
	views := make([]TokenView, 0, len(records))
	for _, record := range records {
		views = append(views, TokenView{
			ID:         record.ID,
			Label:      record.Label,
			Preview:    record.Preview,
			CreatedAt:  record.CreatedAt,
			LastUsedAt: record.LastUsedAt,
		})
	}

	return views
}
