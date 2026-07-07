package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Installation struct {
	ID      int64     `json:"id"`
	Account RepoOwner `json:"account"`
}

func FetchUserInstallations(requestContext context.Context, accessToken string) ([]Installation, error) {
	request, requestError := http.NewRequestWithContext(requestContext, http.MethodGet, UserInstallationsEndpoint, nil)
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
		return nil, fmt.Errorf(InstallationsFetchFailed, response.StatusCode)
	}

	var payload struct {
		Installations []Installation `json:"installations"`
	}
	if decodeError := json.NewDecoder(response.Body).Decode(&payload); decodeError != nil {
		return nil, decodeError
	}

	return payload.Installations, nil
}

func FetchInstallationRepos(requestContext context.Context, accessToken string, installationID int64) ([]Repository, error) {
	endpoint := fmt.Sprintf(InstallationReposEndpointFormat, installationID)

	request, requestError := http.NewRequestWithContext(requestContext, http.MethodGet, endpoint, nil)
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
		return nil, fmt.Errorf(InstallationReposFetchFailed, response.StatusCode)
	}

	var payload struct {
		Repositories []Repository `json:"repositories"`
	}
	if decodeError := json.NewDecoder(response.Body).Decode(&payload); decodeError != nil {
		return nil, decodeError
	}

	return payload.Repositories, nil
}
