package settings

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"git.shi.foo/models"
	"git.shi.foo/repositories/token"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateToken(userID uint, label string) (string, *fiber.Error) {
	if label == "" {
		return "", shortcuts.ServiceError(fiber.StatusBadRequest, LabelRequired)
	}

	secret, generateError := generateSecret()
	if generateError != nil {
		logger.Errorf(LogPrefix, TokenGenerateLog, generateError)
		return "", shortcuts.ServiceError(fiber.StatusInternalServerError, TokenCreateFailed)
	}

	plaintext := TokenPrefix + secret
	record := &models.PersonalAccessToken{
		UserID:    userID,
		Label:     label,
		TokenHash: HashToken(plaintext),
		Preview:   preview(secret),
	}
	if createError := token.Create(record); createError != nil {
		logger.Errorf(LogPrefix, TokenCreateLog, createError)
		return "", shortcuts.ServiceError(fiber.StatusInternalServerError, TokenCreateFailed)
	}

	return plaintext, nil
}

func RevokeToken(userID uint, tokenID uint) *fiber.Error {
	record, findError := token.FindByID(tokenID)
	if findError != nil {
		if errors.Is(findError, gorm.ErrRecordNotFound) {
			return shortcuts.ServiceError(fiber.StatusNotFound, TokenNotFound)
		}
		logger.Errorf(LogPrefix, TokenRevokeLog, findError)
		return shortcuts.ServiceError(fiber.StatusInternalServerError, TokenRevokeFailed)
	}

	if record.UserID != userID {
		return shortcuts.ServiceError(fiber.StatusForbidden, TokenNotYours)
	}

	if deleteError := token.Delete(record); deleteError != nil {
		logger.Errorf(LogPrefix, TokenRevokeLog, deleteError)
		return shortcuts.ServiceError(fiber.StatusInternalServerError, TokenRevokeFailed)
	}

	return nil
}

func HashToken(plaintext string) string {
	sum := sha256.Sum256([]byte(plaintext))
	return hex.EncodeToString(sum[:])
}

func generateSecret() (string, error) {
	buffer := make([]byte, 24)
	if _, readError := rand.Read(buffer); readError != nil {
		return "", readError
	}

	return hex.EncodeToString(buffer), nil
}

func preview(secret string) string {
	return TokenPrefix + secret[:4] + "…" + secret[len(secret)-4:]
}
