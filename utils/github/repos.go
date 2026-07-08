package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type RepoOwner struct {
	Login string `json:"login"`
}

type RepoPermissions struct {
	Admin bool `json:"admin"`
	Push  bool `json:"push"`
	Pull  bool `json:"pull"`
}

type Repository struct {
	ID            int64           `json:"id"`
	Name          string          `json:"name"`
	FullName      string          `json:"full_name"`
	Private       bool            `json:"private"`
	Description   string          `json:"description"`
	DefaultBranch string          `json:"default_branch"`
	Fork          bool            `json:"fork"`
	Archived      bool            `json:"archived"`
	Owner         RepoOwner       `json:"owner"`
	Permissions   RepoPermissions `json:"permissions"`
}

func FetchUserRepos(requestContext context.Context, accessToken string) ([]Repository, error) {
	request, requestError := http.NewRequestWithContext(requestContext, http.MethodGet, UserReposEndpoint, nil)
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
		return nil, fmt.Errorf(UserReposFetchFailed, response.StatusCode)
	}

	var repositories []Repository
	if decodeError := json.NewDecoder(response.Body).Decode(&repositories); decodeError != nil {
		return nil, decodeError
	}

	return repositories, nil
}

func FetchRepo(requestContext context.Context, accessToken string, owner string, name string) (*Repository, error) {
	endpoint := fmt.Sprintf(RepoEndpointFormat, owner, name)

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
		return nil, fmt.Errorf(RepoFetchFailed, response.StatusCode)
	}

	var repository Repository
	if decodeError := json.NewDecoder(response.Body).Decode(&repository); decodeError != nil {
		return nil, decodeError
	}

	return &repository, nil
}
