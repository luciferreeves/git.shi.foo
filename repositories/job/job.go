package job

import (
	"time"

	"git.shi.foo/database"
	"git.shi.foo/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Create(record *models.Job) error {
	return database.DB.Create(record).Error
}

func Update(record *models.Job) error {
	return database.DB.Save(record).Error
}

func FindByID(id uint) (*models.Job, error) {
	var found models.Job
	result := database.DB.First(&found, id)
	return &found, result.Error
}

func ClaimNext(now time.Time) (*models.Job, error) {
	var claimed *models.Job

	transactionError := database.DB.Transaction(func(transaction *gorm.DB) error {
		var candidate models.Job
		selectResult := transaction.
			Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("status = ? AND run_after <= ?", StatusPending, now).
			Order("run_after asc").
			First(&candidate)

		if selectResult.Error != nil {
			return selectResult.Error
		}

		candidate.Status = StatusRunning
		candidate.Attempts = candidate.Attempts + 1
		if saveError := transaction.Save(&candidate).Error; saveError != nil {
			return saveError
		}

		claimed = &candidate
		return nil
	})

	if transactionError != nil {
		if transactionError == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, transactionError
	}

	return claimed, nil
}

func RequeueRunning() error {
	return database.DB.Model(&models.Job{}).
		Where("status = ?", StatusRunning).
		Update("status", StatusPending).Error
}

func FindLatestForRepo(repoID uint, kind string) (*models.Job, error) {
	var found models.Job
	result := database.DB.
		Where("repo_id = ? AND kind = ?", repoID, kind).
		Order("created_at desc").
		First(&found)
	return &found, result.Error
}

func HasOpenForRepo(repoID uint, kind string) (bool, error) {
	var count int64
	result := database.DB.Model(&models.Job{}).
		Where("repo_id = ? AND kind = ? AND status IN ?", repoID, kind, []string{StatusPending, StatusRunning}).
		Count(&count)
	return count > 0, result.Error
}
