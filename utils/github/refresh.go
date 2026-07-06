package github

import (
	"context"

	"golang.org/x/oauth2"
)

func Refresh(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	source := OAuthConfig.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken})
	return source.Token()
}
