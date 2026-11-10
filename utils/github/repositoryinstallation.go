package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func RepositoryInstallationID(requestContext context.Context, owner string, name string) (string, error) {
	token, jwtError := appJWT()
	if jwtError != nil {
		return "", jwtError
	}

	endpoint := fmt.Sprintf(RepoInstallationEndpointFormat, owner, name)
	request, requestError := http.NewRequestWithContext(requestContext, http.MethodGet, endpoint, nil)
	if requestError != nil {
		return "", requestError
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", APIVersion)

	response, responseError := http.DefaultClient.Do(request)
	if responseError != nil {
		return "", responseError
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf(RepoInstallationFetchFailed, response.StatusCode)
	}

	var installation Installation
	if decodeError := json.NewDecoder(response.Body).Decode(&installation); decodeError != nil {
		return "", decodeError
	}
	if installation.ID == 0 {
		return "", errors.New(RepoInstallationMissing)
	}
	return strconv.FormatInt(installation.ID, 10), nil
}
