package settings

import (
	"errors"
	"strings"

	"git.shi.foo/models"
	"git.shi.foo/repositories/key"
	"git.shi.foo/utils/github"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

func AddKey(userID uint, title string, keyText string) *fiber.Error {
	created, storeError := storeKey(userID, title, keyText, key.SourceManual)
	if storeError != nil {
		return storeError
	}

	if !created {
		return shortcuts.ServiceError(fiber.StatusConflict, KeyExists)
	}

	return nil
}

func ImportKeys(userID uint, login string) *fiber.Error {
	fetched, fetchError := github.FetchPublicSSHKeys(login)
	if fetchError != nil {
		logger.Errorf(LogPrefix, KeyImportLog, fetchError)
		return shortcuts.ServiceError(fiber.StatusBadGateway, KeyImportFailed)
	}

	for _, keyText := range fetched {
		if _, storeError := storeKey(userID, "", keyText, key.SourceGitHub); storeError != nil {
			logger.Warnf(LogPrefix, KeyImportSkipLog, keyText)
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

func storeKey(userID uint, title string, keyText string, source string) (bool, *fiber.Error) {
	parsed, comment, _, _, parseError := ssh.ParseAuthorizedKey([]byte(keyText))
	if parseError != nil {
		return false, shortcuts.ServiceError(fiber.StatusBadRequest, InvalidPublicKey)
	}

	fingerprint := ssh.FingerprintSHA256(parsed)

	existing, findError := key.FindByFingerprint(fingerprint)
	if findError == nil && existing != nil {
		return false, nil
	}
	if findError != nil && !errors.Is(findError, gorm.ErrRecordNotFound) {
		logger.Errorf(LogPrefix, KeyAddLog, findError)
		return false, shortcuts.ServiceError(fiber.StatusInternalServerError, KeyAddFailed)
	}

	resolvedTitle := strings.TrimSpace(title)
	if resolvedTitle == "" {
		resolvedTitle = strings.TrimSpace(comment)
	}
	if resolvedTitle == "" {
		resolvedTitle = DefaultKeyTitle
	}

	record := &models.PublicKey{
		UserID:      userID,
		Title:       resolvedTitle,
		KeyType:     parsed.Type(),
		Fingerprint: fingerprint,
		Body:        strings.TrimSpace(string(ssh.MarshalAuthorizedKey(parsed))),
		Source:      source,
	}
	if createError := key.Create(record); createError != nil {
		logger.Errorf(LogPrefix, KeyAddLog, createError)
		return false, shortcuts.ServiceError(fiber.StatusInternalServerError, KeyAddFailed)
	}

	return true, nil
}
