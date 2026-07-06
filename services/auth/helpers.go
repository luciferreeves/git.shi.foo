package auth

import (
	"errors"

	"git.shi.foo/models"
	"git.shi.foo/repositories/credential"
	"git.shi.foo/repositories/user"
	"git.shi.foo/utils/github"

	"gorm.io/gorm"
)

func upsertUser(providerID string, identity *github.Identity) (*models.User, error) {
	existing, findError := user.FindByProviderID(providerID)
	if findError != nil && !errors.Is(findError, gorm.ErrRecordNotFound) {
		return nil, findError
	}

	if findError == nil {
		existing.Login = identity.Login
		existing.Email = identity.Email
		existing.Avatar = identity.AvatarURL
		if updateError := user.Update(existing); updateError != nil {
			return nil, updateError
		}
		return existing, nil
	}

	created := &models.User{
		ProviderID: providerID,
		Login:      identity.Login,
		Email:      identity.Email,
		Avatar:     identity.AvatarURL,
	}
	if createError := user.Create(created); createError != nil {
		return nil, createError
	}

	return created, nil
}

func storeRefreshToken(userID uint, refreshToken string) error {
	ciphertext, encryptError := github.Encrypt(refreshToken)
	if encryptError != nil {
		return encryptError
	}

	return credential.Upsert(userID, ciphertext)
}
