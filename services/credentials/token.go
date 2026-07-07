package credentials

import (
	"context"
	"fmt"
	"time"

	"git.shi.foo/repositories/cache"
	"git.shi.foo/repositories/credential"
	"git.shi.foo/utils/github"
	"git.shi.foo/utils/logger"
)

func AccessTokenForUser(requestContext context.Context, userID uint) (string, error) {
	cacheKey := fmt.Sprintf(AccessCacheKeyFormat, userID)

	if entry, cacheError := cache.Get(cacheKey); cacheError == nil {
		return string(entry.Value), nil
	}

	stored, findError := credential.FindByUserID(userID)
	if findError != nil {
		logger.Errorf(LogPrefix, CredentialMissingLog, findError)
		return "", findError
	}

	refreshToken, decryptError := github.Decrypt(stored.RefreshToken)
	if decryptError != nil {
		logger.Errorf(LogPrefix, DecryptLog, decryptError)
		return "", decryptError
	}

	minted, refreshError := github.Refresh(requestContext, refreshToken)
	if refreshError != nil {
		logger.Errorf(LogPrefix, RefreshLog, refreshError)
		return "", refreshError
	}

	rotateRefreshToken(userID, minted.RefreshToken)
	cacheAccessToken(cacheKey, minted.AccessToken, minted.Expiry)

	return minted.AccessToken, nil
}

func rotateRefreshToken(userID uint, refreshToken string) {
	if refreshToken == "" {
		return
	}

	ciphertext, encryptError := github.Encrypt(refreshToken)
	if encryptError != nil {
		logger.Errorf(LogPrefix, RotateLog, encryptError)
		return
	}

	if upsertError := credential.Upsert(userID, ciphertext); upsertError != nil {
		logger.Errorf(LogPrefix, RotateLog, upsertError)
	}
}

func cacheAccessToken(cacheKey string, accessToken string, expiry time.Time) {
	freshUntil := expiry.Add(-AccessSafetyMargin)
	if expiry.IsZero() {
		freshUntil = time.Now().Add(AccessFallbackTTL)
	}

	if cacheError := cache.Set(cacheKey, []byte(accessToken), freshUntil); cacheError != nil {
		logger.Errorf(LogPrefix, CacheLog, cacheError)
	}
}
