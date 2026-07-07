package github

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"time"

	"git.shi.foo/config"
)

type installationToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func appJWT() (string, error) {
	block, _ := pem.Decode([]byte(config.GitHub.AppPrivateKey))
	if block == nil {
		return "", errors.New(PrivateKeyInvalid)
	}

	privateKey, parseError := x509.ParsePKCS1PrivateKey(block.Bytes)
	if parseError != nil {
		return "", parseError
	}

	now := time.Now()
	header := base64URL([]byte(`{"alg":"RS256","typ":"JWT"}`))
	claims := base64URL([]byte(fmt.Sprintf(`{"iat":%d,"exp":%d,"iss":"%s"}`,
		now.Add(-AppJWTSkew).Unix(),
		now.Add(AppJWTLifetime).Unix(),
		config.GitHub.AppID,
	)))

	signingInput := header + "." + claims
	digest := sha256.Sum256([]byte(signingInput))

	signature, signError := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, digest[:])
	if signError != nil {
		return "", signError
	}

	return signingInput + "." + base64URL(signature), nil
}

func InstallationToken(installationID string) (string, time.Time, error) {
	token, jwtError := appJWT()
	if jwtError != nil {
		return "", time.Time{}, jwtError
	}

	endpoint := fmt.Sprintf("%s/app/installations/%s/access_tokens", APIBase, installationID)

	request, requestError := http.NewRequest(http.MethodPost, endpoint, nil)
	if requestError != nil {
		return "", time.Time{}, requestError
	}

	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", APIVersion)

	response, responseError := http.DefaultClient.Do(request)
	if responseError != nil {
		return "", time.Time{}, responseError
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return "", time.Time{}, fmt.Errorf(InstallationTokenFailed, response.StatusCode)
	}

	var minted installationToken
	if decodeError := json.NewDecoder(response.Body).Decode(&minted); decodeError != nil {
		return "", time.Time{}, decodeError
	}

	return minted.Token, minted.ExpiresAt, nil
}

func base64URL(input []byte) string {
	return base64.RawURLEncoding.EncodeToString(input)
}
