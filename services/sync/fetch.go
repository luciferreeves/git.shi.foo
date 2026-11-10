package sync

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"git.shi.foo/git"
	"git.shi.foo/jobs"
	"git.shi.foo/models"
	"git.shi.foo/repositories/cache"
	"git.shi.foo/repositories/repo"
	"git.shi.foo/services/credentials"
	"git.shi.foo/services/repos"
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

	var token string
	var tokenError error
	if payload.InstallationID == "" {
		token, tokenError = credentials.RepositoryAccessToken(runContext, repoRecord.Owner, repoRecord.Name)
	} else {
		token, tokenError = credentials.InstallationAccessToken(runContext, payload.InstallationID)
	}
	if tokenError != nil {
		return tokenError
	}

	metadata, metadataError := github.FetchRepo(runContext, token, repoRecord.Owner, repoRecord.Name)
	if metadataError != nil {
		return metadataError
	}
	if metadata.ID == 0 || metadata.ID != repoRecord.GitHubID {
		return fmt.Errorf(RepositoryIdentityMismatch, repoRecord.Owner, repoRecord.Name)
	}
	if repoRecord.Description != metadata.Description || repoRecord.Private != metadata.Private || (metadata.DefaultBranch != "" && repoRecord.DefaultBranch != metadata.DefaultBranch) {
		repoRecord.Description = metadata.Description
		repoRecord.Private = metadata.Private
		if metadata.DefaultBranch != "" {
			repoRecord.DefaultBranch = metadata.DefaultBranch
		}
		if updateError := repo.UpdateMetadata(repoRecord); updateError != nil {
			return updateError
		}
	}

	if fetchError := git.Fetch(repoRecord.Owner, repoRecord.Name, token, func(phase string, percent int) {
		report.Progress(phase, percent)
	}); fetchError != nil {
		return fetchError
	}

	return cache.Delete(fmt.Sprintf(repos.InsightsCacheKey, repoRecord.Owner, repoRecord.Name))
}
