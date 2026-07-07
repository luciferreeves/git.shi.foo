package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type emailEntry struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func FetchPrimaryEmail(requestContext context.Context, token *oauth2.Token) (string, error) {
	request, requestError := http.NewRequestWithContext(requestContext, http.MethodGet, UserEmailsEndpoint, nil)
	if requestError != nil {
		return "", requestError
	}

	request.Header.Set("Authorization", "Bearer "+token.AccessToken)
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", APIVersion)

	response, responseError := http.DefaultClient.Do(request)
	if responseError != nil {
		return "", responseError
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf(UserEmailsFetchFailed, response.StatusCode)
	}

	var entries []emailEntry
	if decodeError := json.NewDecoder(response.Body).Decode(&entries); decodeError != nil {
		return "", decodeError
	}

	fallback := ""
	for _, entry := range entries {
		if entry.Primary && entry.Verified {
			return entry.Email, nil
		}
		if fallback == "" && entry.Verified {
			fallback = entry.Email
		}
	}

	return fallback, nil
}
