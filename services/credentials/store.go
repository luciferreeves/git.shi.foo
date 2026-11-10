package credentials

import (
	"fmt"

	"git.shi.foo/repositories/credential"
	"git.shi.foo/utils/github"

	"golang.org/x/oauth2"
)

func StoreTokenPair(userID uint, token *oauth2.Token) error {
	tokenLock := tokenLockForUser(userID)
	tokenLock.Lock()
	defer tokenLock.Unlock()

	return persistTokenPair(userID, token)
}

func persistTokenPair(userID uint, token *oauth2.Token) error {
	ciphertext, encryptError := github.Encrypt(token.RefreshToken)
	if encryptError != nil {
		return encryptError
	}
	if upsertError := credential.Upsert(userID, ciphertext); upsertError != nil {
		return upsertError
	}

	return cacheAccessToken(fmt.Sprintf(AccessCacheKeyFormat, userID), token.AccessToken, token.Expiry)
}
