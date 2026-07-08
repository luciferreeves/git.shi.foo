package sync

import (
	"context"
	"encoding/json"
	"errors"

	"git.shi.foo/git"
	"git.shi.foo/jobs"
	"git.shi.foo/models"
	"git.shi.foo/repositories/repo"
	"git.shi.foo/utils/github"
)

func RunFetch(runContext context.Context, record *models.Job, report jobs.Reporter) error {
	if record.RepoID == nil {
		return errors.New(FetchMissingRepo)
	}

	repoRecord, findError := repo.FindByID(*record.RepoID)
	if findError != nil {
		return findError
	}

	var payload FetchPayload
	if len(record.Payload) > 0 {
		if unmarshalError := json.Unmarshal(record.Payload, &payload); unmarshalError != nil {
			return unmarshalError
		}
	}

	if payload.InstallationID == "" {
		return errors.New(MissingInstallation)
	}

	token, _, tokenError := github.InstallationToken(payload.InstallationID)
	if tokenError != nil {
		return tokenError
	}

	return git.Fetch(repoRecord.Owner, repoRecord.Name, token, func(phase string, percent int) {
		report.Progress(phase, percent)
	})
}
