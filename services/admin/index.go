package admin

import (
	"git.shi.foo/account"
	"git.shi.foo/models"
	"git.shi.foo/repositories/invitation"
	"git.shi.foo/repositories/user"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func GetIndexData(currentUser *account.Response) (*IndexContext, *fiber.Error) {
	if guardError := EnsureAdmin(currentUser); guardError != nil {
		return nil, guardError
	}

	users, usersError := user.All()
	if usersError != nil {
		logger.Errorf(LogPrefix, LoadLog, usersError)
		return nil, shortcuts.ServiceError(fiber.StatusInternalServerError, LoadFailed)
	}

	invitations, invitationsError := invitation.All()
	if invitationsError != nil {
		logger.Errorf(LogPrefix, LoadLog, invitationsError)
		return nil, shortcuts.ServiceError(fiber.StatusInternalServerError, LoadFailed)
	}

	return &IndexContext{
		Title:       AdminTitle,
		Users:       toUserViews(users),
		Invitations: toInvitationViews(invitations),
	}, nil
}

func toUserViews(records []models.User) []UserView {
	views := make([]UserView, 0, len(records))
	for _, record := range records {
		views = append(views, UserView{
			Login:     record.Login,
			Email:     record.Email,
			Admin:     record.Admin,
			Enabled:   record.Enabled,
			CreatedAt: record.CreatedAt,
		})
	}

	return views
}

func toInvitationViews(records []models.Invitation) []InvitationView {
	views := make([]InvitationView, 0, len(records))
	for _, record := range records {
		views = append(views, InvitationView{
			ID:        record.ID,
			Email:     record.Email,
			Username:  record.Username,
			Status:    record.Status,
			CreatedAt: record.CreatedAt,
		})
	}

	return views
}
