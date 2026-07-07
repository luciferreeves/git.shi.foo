package settings

import (
	"git.shi.foo/account"
	"git.shi.foo/models"
	"git.shi.foo/repositories/key"
	"git.shi.foo/repositories/token"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func GetIndexData(currentUser *account.Response) (*IndexContext, *fiber.Error) {
	if guardError := EnsureUser(currentUser); guardError != nil {
		return nil, guardError
	}

	tokens, tokenListError := token.ListByUser(currentUser.ID)
	if tokenListError != nil {
		logger.Errorf(LogPrefix, TokenListLog, tokenListError)
		return nil, shortcuts.ServiceError(fiber.StatusInternalServerError, TokenListFailed)
	}

	keys, keyListError := key.ListByUser(currentUser.ID)
	if keyListError != nil {
		logger.Errorf(LogPrefix, KeyListLog, keyListError)
		return nil, shortcuts.ServiceError(fiber.StatusInternalServerError, KeyListFailed)
	}

	return &IndexContext{
		Title:  SettingsTitle,
		Tokens: toTokenViews(tokens),
		Keys:   toKeyViews(keys),
	}, nil
}

func toKeyViews(records []models.PublicKey) []KeyView {
	views := make([]KeyView, 0, len(records))
	for _, record := range records {
		views = append(views, KeyView{
			ID:          record.ID,
			Title:       record.Title,
			KeyType:     record.KeyType,
			Fingerprint: record.Fingerprint,
			Source:      record.Source,
			CreatedAt:   record.CreatedAt,
		})
	}

	return views
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
