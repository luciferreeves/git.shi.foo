package repos

import (
	jobrepo "git.shi.foo/repositories/job"
	"git.shi.foo/repositories/repo"
	"git.shi.foo/utils/logger"
)

func ReconcileImports() {
	records, listError := repo.ListByStatus(repo.StatusImporting)
	if listError != nil {
		logger.Errorf(LogPrefix, ReconcileLog, listError)
		return
	}

	for index := range records {
		record := records[index]

		open, checkError := jobrepo.HasOpenForRepo(record.ID, jobrepo.KindImport)
		if checkError != nil {
			logger.Errorf(LogPrefix, ReconcileLog, checkError)
			continue
		}
		if open {
			continue
		}

		if enqueueError := enqueueImport(record.ID); enqueueError != nil {
			logger.Errorf(LogPrefix, EnqueueLog, enqueueError)
		}
	}
}
