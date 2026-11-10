package sync

import "git.shi.foo/repositories/repo"

func Reconcile() error {
	records, listError := repo.ListActive()
	if listError != nil {
		return listError
	}
	for _, record := range records {
		if enqueueError := enqueueFetch(record.ID, 0); enqueueError != nil {
			return enqueueError
		}
	}
	return nil
}
