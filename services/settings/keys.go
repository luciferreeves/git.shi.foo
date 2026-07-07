package settings

import (
	"context"
	"errors"
	"strings"

	"git.shi.foo/models"
	"git.shi.foo/repositories/key"
	"git.shi.foo/services/credentials"
	"git.shi.foo/utils/github"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

func AddKey(userID uint, title string, keyText string) *fiber.Error {
	created, storeError := storeKey(userID, title, keyText, key.SourceManual, 0)
	if storeError != nil {
		return storeError
	}

	if !created {
		return shortcuts.ServiceError(fiber.StatusConflict, KeyExists)
	}

	return nil
}

func ImportKeys(requestContext context.Context, userID uint) *fiber.Error {
	accessToken, tokenError := credentials.AccessTokenForUser(requestContext, userID)
	if tokenError != nil {
		logger.Errorf(LogPrefix, KeyImportLog, tokenError)
		return shortcuts.ServiceError(fiber.StatusBadGateway, KeyImportFailed)
	}

	fetched, fetchError := github.FetchUserSSHKeys(requestContext, accessToken)
	if fetchError != nil {
		logger.Errorf(LogPrefix, KeyImportLog, fetchError)
		return shortcuts.ServiceError(fiber.StatusBadGateway, KeyImportFailed)
	}

	for _, item := range fetched {
		if _, storeError := storeKey(userID, item.Title, item.Key, key.SourceGitHub, item.ID); storeError != nil {
			logger.Warnf(LogPrefix, KeyImportSkipLog, item.Key)
		}
	}

	return nil
}

func RemoveKey(userID uint, keyID uint) *fiber.Error {
	record, findError := key.FindByID(keyID)
	if findError != nil {
		if errors.Is(findError, gorm.ErrRecordNotFound) {
			return shortcuts.ServiceError(fiber.StatusNotFound, KeyNotFound)
		}
		logger.Errorf(LogPrefix, KeyRemoveLog, findError)
		return shortcuts.ServiceError(fiber.StatusInternalServerError, KeyRemoveFailed)
	}

	if record.UserID != userID {
		return shortcuts.ServiceError(fiber.StatusForbidden, KeyNotYours)
	}

	if deleteError := key.Delete(record); deleteError != nil {
		logger.Errorf(LogPrefix, KeyRemoveLog, deleteError)
		return shortcuts.ServiceError(fiber.StatusInternalServerError, KeyRemoveFailed)
	}

	return nil
}

func storeKey(userID uint, title string, keyText string, source string, githubID int64) (bool, *fiber.Error) {
	parsed, comment, _, _, parseError := ssh.ParseAuthorizedKey([]byte(keyText))
	if parseError != nil {
		return false, shortcuts.ServiceError(fiber.StatusBadRequest, InvalidPublicKey)
	}

	fingerprint := ssh.FingerprintSHA256(parsed)

	resolvedTitle := strings.TrimSpace(title)
	if resolvedTitle == "" {
		resolvedTitle = strings.TrimSpace(comment)
	}
	if resolvedTitle == "" {
		resolvedTitle = DefaultKeyTitle
	}

	existing, findError := key.FindByFingerprint(fingerprint)
	if findError == nil && existing != nil {
		if existing.UserID == userID && title != "" && existing.Title != resolvedTitle {
			existing.Title = resolvedTitle
			existing.Source = source
			existing.GitHubID = githubID
			if updateError := key.Update(existing); updateError != nil {
				logger.Errorf(LogPrefix, KeyAddLog, updateError)
				return false, shortcuts.ServiceError(fiber.StatusInternalServerError, KeyAddFailed)
			}
		}
		return false, nil
	}
	if findError != nil && !errors.Is(findError, gorm.ErrRecordNotFound) {
		logger.Errorf(LogPrefix, KeyAddLog, findError)
		return false, shortcuts.ServiceError(fiber.StatusInternalServerError, KeyAddFailed)
	}

	record := &models.PublicKey{
		UserID:      userID,
		Title:       resolvedTitle,
		KeyType:     parsed.Type(),
		Fingerprint: fingerprint,
		Body:        strings.TrimSpace(string(ssh.MarshalAuthorizedKey(parsed))),
		Source:      source,
		GitHubID:    githubID,
	}
	if createError := key.Create(record); createError != nil {
		logger.Errorf(LogPrefix, KeyAddLog, createError)
		return false, shortcuts.ServiceError(fiber.StatusInternalServerError, KeyAddFailed)
	}

	return true, nil
}
