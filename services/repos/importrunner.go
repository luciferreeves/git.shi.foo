package repos

import (
	"context"
	"errors"

	"git.shi.foo/git"
	"git.shi.foo/jobs"
	"git.shi.foo/models"
	"git.shi.foo/repositories/repo"
	"git.shi.foo/services/credentials"
	"git.shi.foo/utils/logger"
)

func RunImport(runContext context.Context, record *models.Job, report jobs.Reporter) error {
	if record.RepoID == nil {
		return errors.New(ImportMissingRepo)
	}

	repoRecord, findError := repo.FindByID(*record.RepoID)
	if findError != nil {
		return findError
	}

	if mirrorError := mirrorRepo(runContext, repoRecord, report); mirrorError != nil {
		if report.Final() {
			markRepo(repoRecord, repo.StatusFailed)
		}
		return mirrorError
	}

	repoRecord.Status = repo.StatusActive
	return repo.Update(repoRecord)
}

func mirrorRepo(runContext context.Context, repoRecord *models.Repo, report jobs.Reporter) error {
	accessToken, tokenError := credentials.AccessTokenForUser(runContext, repoRecord.ImportedBy)
	if tokenError != nil {
		return tokenError
	}

	return git.MirrorWithProgress(repoRecord.Owner, repoRecord.Name, accessToken, func(phase string, percent int) {
		report.Progress(phase, percent)
	})
}

func markRepo(repoRecord *models.Repo, status string) {
	repoRecord.Status = status
	if updateError := repo.Update(repoRecord); updateError != nil {
		logger.Errorf(LogPrefix, StatusUpdateLog, updateError)
	}
}
