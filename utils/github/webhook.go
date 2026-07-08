package github

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

const SignaturePrefix = "sha256="

type WebhookInstallation struct {
	ID int64 `json:"id"`
}

type WebhookRepository struct {
	Name  string    `json:"name"`
	Owner RepoOwner `json:"owner"`
}

type WebhookPayload struct {
	Repository   WebhookRepository   `json:"repository"`
	Installation WebhookInstallation `json:"installation"`
}

func VerifySignature(payload []byte, signature string, secret string) bool {
	if secret == "" || signature == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expected := SignaturePrefix + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(signature))
}
