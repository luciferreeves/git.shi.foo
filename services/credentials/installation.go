package credentials

import (
	"context"
	"fmt"
	"sync"

	"git.shi.foo/repositories/cache"
	"git.shi.foo/utils/github"
)

var installationTokenLock sync.Mutex

func InstallationAccessToken(requestContext context.Context, installationID string) (string, error) {
	installationTokenLock.Lock()
	defer installationTokenLock.Unlock()

	cacheKey := fmt.Sprintf(InstallationCacheKeyFormat, installationID)
	if entry, cacheError := cache.Get(cacheKey); cacheError == nil {
		return string(entry.Value), nil
	}

	token, expiry, tokenError := github.InstallationTokenWithContext(requestContext, installationID)
	if tokenError != nil {
		return "", tokenError
	}
	if storeError := cache.Set(cacheKey, []byte(token), expiry.Add(-AccessSafetyMargin)); storeError != nil {
		return "", storeError
	}
	return token, nil
}
