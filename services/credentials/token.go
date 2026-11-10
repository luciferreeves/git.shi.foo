package credentials

import (
	"context"
	"errors"
	"fmt"
	"time"

	"git.shi.foo/repositories/cache"
	"git.shi.foo/repositories/credential"
	"git.shi.foo/utils/github"
	"git.shi.foo/utils/logger"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

func AccessTokenForUser(requestContext context.Context, userID uint) (string, error) {
	tokenLock := tokenLockForUser(userID)
	tokenLock.Lock()
	defer tokenLock.Unlock()

	cacheKey := fmt.Sprintf(AccessCacheKeyFormat, userID)
	if entry, cacheError := cache.Get(cacheKey); cacheError == nil {
		return string(entry.Value), nil
	}

	return mintAccessToken(requestContext, userID)
}

func RefreshAccessTokenForUser(requestContext context.Context, userID uint) (string, error) {
	tokenLock := tokenLockForUser(userID)
	tokenLock.Lock()
	defer tokenLock.Unlock()

	return mintAccessToken(requestContext, userID)
}

func mintAccessToken(requestContext context.Context, userID uint) (string, error) {
	stored, findError := credential.FindByUserID(userID)
	if findError != nil {
		logger.Errorf(LogPrefix, CredentialMissingLog, userID, findError)
		return "", findError
	}
	if stored.ReauthorizationRequired {
		return "", fiber.NewError(fiber.StatusUnauthorized, ReauthorizationRequired)
	}

	refreshToken, decryptError := github.Decrypt(stored.RefreshToken)
	if decryptError != nil {
		logger.Errorf(LogPrefix, DecryptLog, userID, decryptError)
		return "", decryptError
	}

	minted, refreshError := github.Refresh(requestContext, refreshToken)
	if refreshError != nil {
		var retrievalError *oauth2.RetrieveError
		if errors.As(refreshError, &retrievalError) && (retrievalError.ErrorCode == InvalidRefreshTokenCode || retrievalError.ErrorCode == InvalidGrantCode) {
			if markError := credential.RequireReauthorization(userID); markError != nil {
				logger.Errorf(LogPrefix, ReauthorizationStoreLog, userID, markError)
				return "", markError
			}
			if deleteError := cache.Delete(fmt.Sprintf(AccessCacheKeyFormat, userID)); deleteError != nil {
				logger.Errorf(LogPrefix, CacheLog, userID, deleteError)
				return "", deleteError
			}
			logger.Warnf(LogPrefix, ReauthorizationLog, userID)
			return "", fiber.NewError(fiber.StatusUnauthorized, ReauthorizationRequired)
		}
		logger.Errorf(LogPrefix, RefreshLog, userID, refreshError)
		return "", refreshError
	}

	if storeError := persistTokenPair(userID, minted); storeError != nil {
		logger.Errorf(LogPrefix, RotateLog, userID, storeError)
		return "", storeError
	}

	return minted.AccessToken, nil
}

func cacheAccessToken(cacheKey string, accessToken string, expiry time.Time) error {
	freshUntil := expiry.Add(-AccessSafetyMargin)
	if expiry.IsZero() {
		freshUntil = time.Now().Add(AccessFallbackTTL)
	}

	return cache.Set(cacheKey, []byte(accessToken), freshUntil)
}
