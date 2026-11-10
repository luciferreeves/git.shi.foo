package credentials

import (
	"context"
	"fmt"
	"time"

	"git.shi.foo/repositories/cache"
	"git.shi.foo/utils/github"
)

func RepositoryAccessToken(requestContext context.Context, owner string, name string) (string, error) {
	cacheKey := fmt.Sprintf(RepositoryInstallationCacheKeyFormat, owner, name)
	var installationID string
	if entry, cacheError := cache.Get(cacheKey); cacheError == nil {
		installationID = string(entry.Value)
	} else {
		resolvedID, resolveError := github.RepositoryInstallationID(requestContext, owner, name)
		if resolveError != nil {
			return "", resolveError
		}
		installationID = resolvedID
		if storeError := cache.Set(cacheKey, []byte(installationID), time.Now().Add(RepositoryInstallationTTL)); storeError != nil {
			return "", storeError
		}
	}
	return InstallationAccessToken(requestContext, installationID)
}
