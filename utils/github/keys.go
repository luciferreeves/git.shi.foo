package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type SSHKey struct {
	ID    int64  `json:"id"`
	Key   string `json:"key"`
	Title string `json:"title"`
}

func FetchUserSSHKeys(requestContext context.Context, accessToken string) ([]SSHKey, error) {
	request, requestError := http.NewRequestWithContext(requestContext, http.MethodGet, UserKeysEndpoint, nil)
	if requestError != nil {
		return nil, requestError
	}

	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", APIVersion)

	response, responseError := http.DefaultClient.Do(request)
	if responseError != nil {
		return nil, responseError
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(UserKeysFetchFailed, response.StatusCode)
	}

	var keys []SSHKey
	if decodeError := json.NewDecoder(response.Body).Decode(&keys); decodeError != nil {
		return nil, decodeError
	}

	return keys, nil
}
