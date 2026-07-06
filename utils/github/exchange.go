package github

import (
	"context"
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"
)

type Identity struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

func Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return OAuthConfig.Exchange(ctx, code)
}

func FetchIdentity(ctx context.Context, token *oauth2.Token) (*Identity, error) {
	request, requestError := http.NewRequestWithContext(ctx, http.MethodGet, UserEndpoint, nil)
	if requestError != nil {
		return nil, requestError
	}

	request.Header.Set("Authorization", "Bearer "+token.AccessToken)
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", APIVersion)

	response, responseError := http.DefaultClient.Do(request)
	if responseError != nil {
		return nil, responseError
	}
	defer response.Body.Close()

	var identity Identity
	if decodeError := json.NewDecoder(response.Body).Decode(&identity); decodeError != nil {
		return nil, decodeError
	}

	return &identity, nil
}
