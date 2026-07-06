package invites

import (
	"fmt"

	"git.shi.foo/config"
	"git.shi.foo/models"
	"git.shi.foo/repositories/invitation"
	"git.shi.foo/utils/github"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/mail"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func CreateInvite(inviterID uint, email string, username string) *fiber.Error {
	if email == "" || username == "" {
		return shortcuts.ServiceError(fiber.StatusBadRequest, MissingFields)
	}

	token, tokenError := github.GenerateState()
	if tokenError != nil {
		logger.Errorf(LogPrefix, TokenGenerationLog, tokenError)
		return shortcuts.ServiceError(fiber.StatusInternalServerError, InviteCreateFailed)
	}

	record := &models.Invitation{
		Email:     email,
		Username:  username,
		Token:     token,
		Status:    invitation.StatusPending,
		InvitedBy: inviterID,
	}
	if createError := invitation.Create(record); createError != nil {
		logger.Errorf(LogPrefix, InviteCreateLog, createError)
		return shortcuts.ServiceError(fiber.StatusInternalServerError, InviteCreateFailed)
	}

	acceptURL := fmt.Sprintf("%s/invites/accept?token=%s", config.Server.PublicURL, token)
	body := fmt.Sprintf(InviteEmailBody, username, acceptURL)
	if sendError := mail.Send(email, InviteEmailSubject, body); sendError != nil {
		logger.Errorf(LogPrefix, InviteEmailLog, sendError)
		return shortcuts.ServiceError(fiber.StatusInternalServerError, InviteEmailFailed)
	}

	return nil
}
