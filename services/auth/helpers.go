package auth

import (
	"errors"

	"git.shi.foo/models"
	"git.shi.foo/repositories/credential"
	"git.shi.foo/repositories/invitation"
	"git.shi.foo/repositories/user"
	"git.shi.foo/utils/github"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func admitUser(providerID string, identity *github.Identity) (*models.User, *fiber.Error) {
	existing, findError := user.FindByProviderID(providerID)
	if findError == nil {
		existing.Login = identity.Login
		existing.Email = identity.Email
		existing.Avatar = identity.AvatarURL
		if updateError := user.Update(existing); updateError != nil {
			logger.Errorf(LogPrefix, UserUpsertLog, updateError)
			return nil, shortcuts.ServiceError(fiber.StatusInternalServerError, UserUpsertFailed)
		}
		return existing, nil
	}

	if !errors.Is(findError, gorm.ErrRecordNotFound) {
		logger.Errorf(LogPrefix, UserUpsertLog, findError)
		return nil, shortcuts.ServiceError(fiber.StatusInternalServerError, UserUpsertFailed)
	}

	total, countError := user.Count()
	if countError != nil {
		logger.Errorf(LogPrefix, AccessCheckLog, countError)
		return nil, shortcuts.ServiceError(fiber.StatusInternalServerError, AccessCheckFailed)
	}

	if total == 0 {
		return createMember(providerID, identity, true)
	}

	invite, inviteError := invitation.FindByUsernameAndStatus(identity.Login, invitation.StatusAccepted)
	if inviteError != nil {
		logger.Warnf(LogPrefix, InvitationRequiredLog, identity.Login)
		return nil, shortcuts.ServiceError(fiber.StatusForbidden, InvitationRequired)
	}

	member, admitError := createMember(providerID, identity, false)
	if admitError != nil {
		return nil, admitError
	}

	invite.Status = invitation.StatusConsumed
	if consumeError := invitation.Update(invite); consumeError != nil {
		logger.Errorf(LogPrefix, InviteConsumeLog, consumeError)
	}

	return member, nil
}

func createMember(providerID string, identity *github.Identity, admin bool) (*models.User, *fiber.Error) {
	created := &models.User{
		ProviderID: providerID,
		Login:      identity.Login,
		Email:      identity.Email,
		Avatar:     identity.AvatarURL,
		Admin:      admin,
		Enabled:    true,
	}
	if createError := user.Create(created); createError != nil {
		logger.Errorf(LogPrefix, UserUpsertLog, createError)
		return nil, shortcuts.ServiceError(fiber.StatusInternalServerError, UserUpsertFailed)
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
