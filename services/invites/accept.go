package invites

import (
	"errors"

	"git.shi.foo/repositories/invitation"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AcceptInvite(token string) *fiber.Error {
	if token == "" {
		return shortcuts.ServiceError(fiber.StatusBadRequest, MissingToken)
	}

	record, findError := invitation.FindByToken(token)
	if findError != nil {
		if errors.Is(findError, gorm.ErrRecordNotFound) {
			return shortcuts.ServiceError(fiber.StatusNotFound, InviteNotFound)
		}
		logger.Errorf(LogPrefix, InviteLookupLog, findError)
		return shortcuts.ServiceError(fiber.StatusInternalServerError, InviteAcceptFailed)
	}

	if record.Status == invitation.StatusConsumed {
		return shortcuts.ServiceError(fiber.StatusConflict, InviteAlreadyUsed)
	}

	if record.Status == invitation.StatusPending {
		record.Status = invitation.StatusAccepted
		if updateError := invitation.Update(record); updateError != nil {
			logger.Errorf(LogPrefix, InviteAcceptLog, updateError)
			return shortcuts.ServiceError(fiber.StatusInternalServerError, InviteAcceptFailed)
		}
	}

	return nil
}
