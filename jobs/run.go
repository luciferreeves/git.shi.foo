package jobs

import (
	"context"
	"fmt"
	"time"

	"git.shi.foo/models"
	jobrepo "git.shi.foo/repositories/job"
	"git.shi.foo/utils/logger"
)

func process(runContext context.Context, record *models.Job) {
	runner, found := lookup(record.Kind)
	if !found {
		terminate(record, jobrepo.StatusFailed, fmt.Sprintf(UnknownKind, record.Kind))
		return
	}

	report := newReporter(record)
	runError := invoke(runContext, runner, record, report)

	if runError == nil {
		record.Percent = 100
		terminate(record, jobrepo.StatusSucceeded, "")
		return
	}

	if record.Attempts >= record.MaxAttempts {
		terminate(record, jobrepo.StatusFailed, runError.Error())
		return
	}

	requeue(record, runError.Error())
}

func invoke(runContext context.Context, runner Runner, record *models.Job, report Reporter) (runError error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			logger.Errorf(LogPrefix, PanicLog, record.ID, recovered)
			runError = fmt.Errorf(JobPanicked, recovered)
		}
	}()

	return runner(runContext, record, report)
}

func terminate(record *models.Job, status string, message string) {
	record.Status = status
	record.Error = message
	if updateError := jobrepo.Update(record); updateError != nil {
		logger.Errorf(LogPrefix, PersistFailed, updateError)
	}
	publish(record, status)
}

func requeue(record *models.Job, message string) {
	record.Status = jobrepo.StatusPending
	record.Error = message
	record.RunAfter = time.Now().Add(time.Duration(record.Attempts) * RetryBackoffBase)
	if updateError := jobrepo.Update(record); updateError != nil {
		logger.Errorf(LogPrefix, PersistFailed, updateError)
	}
	Signal()
}
