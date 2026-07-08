package sync

import (
	"encoding/json"
	"strconv"
	"time"

	"git.shi.foo/config"
	"git.shi.foo/jobs"
	"git.shi.foo/models"
	jobrepo "git.shi.foo/repositories/job"
	"git.shi.foo/repositories/repo"
	"git.shi.foo/utils/github"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
)

func Ingest(eventType string, signature string, body []byte) *fiber.Error {
	if !github.VerifySignature(body, signature, config.GitHub.WebhookSecret) {
		return shortcuts.ServiceError(fiber.StatusUnauthorized, InvalidSignature)
	}

	switch eventType {
	case EventPush, EventCreate, EventDelete:
		handleRefChange(body)
	}

	return nil
}

func handleRefChange(body []byte) {
	var payload github.WebhookPayload
	if unmarshalError := json.Unmarshal(body, &payload); unmarshalError != nil {
		return
	}

	owner := payload.Repository.Owner.Login
	name := payload.Repository.Name
	if owner == "" || name == "" {
		return
	}

	record, findError := repo.FindByOwnerName(owner, name)
	if findError != nil {
		return
	}

	if enqueueError := enqueueFetch(record.ID, payload.Installation.ID); enqueueError != nil {
		logger.Errorf(LogPrefix, EnqueueLog, enqueueError)
	}
}

func enqueueFetch(repoID uint, installationID int64) error {
	open, checkError := jobrepo.HasOpenForRepo(repoID, jobrepo.KindFetch)
	if checkError != nil {
		return checkError
	}
	if open {
		return nil
	}

	payload, marshalError := json.Marshal(FetchPayload{InstallationID: strconv.FormatInt(installationID, 10)})
	if marshalError != nil {
		return marshalError
	}

	jobRecord := &models.Job{
		Kind:        jobrepo.KindFetch,
		RepoID:      &repoID,
		Status:      jobrepo.StatusPending,
		MaxAttempts: FetchMaxAttempts,
		RunAfter:    time.Now(),
		Payload:     datatypes.JSON(payload),
	}
	if createError := jobrepo.Create(jobRecord); createError != nil {
		return createError
	}

	jobs.Signal()
	return nil
}
