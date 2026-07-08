package jobs

import (
	"fmt"

	"git.shi.foo/events"
	"git.shi.foo/models"
	jobrepo "git.shi.foo/repositories/job"
	"git.shi.foo/utils/logger"
)

type progressReporter struct {
	record      *models.Job
	lastPhase   string
	lastPercent int
	final       bool
}

func newReporter(record *models.Job) *progressReporter {
	return &progressReporter{
		record:      record,
		lastPhase:   record.Phase,
		lastPercent: record.Percent,
		final:       record.Attempts >= record.MaxAttempts,
	}
}

func (self *progressReporter) Progress(phase string, percent int) {
	self.record.Phase = phase
	self.record.Percent = percent
	publish(self.record, jobrepo.StatusRunning)

	if phase != self.lastPhase || percent-self.lastPercent >= PersistThreshold || percent >= 100 {
		self.lastPhase = phase
		self.lastPercent = percent
		if updateError := jobrepo.Update(self.record); updateError != nil {
			logger.Errorf(LogPrefix, PersistFailed, updateError)
		}
	}
}

func (self *progressReporter) Final() bool {
	return self.final
}

func Topic(jobID uint) string {
	return fmt.Sprintf(TopicFormat, jobID)
}

func publish(record *models.Job, status string) {
	events.Default.Publish(Topic(record.ID), events.Event{
		Name: EventProgress,
		Data: Progress{
			JobID:   record.ID,
			Status:  status,
			Phase:   record.Phase,
			Percent: record.Percent,
			Done:    IsTerminal(status),
		},
	})
}

func IsTerminal(status string) bool {
	return status == jobrepo.StatusSucceeded || status == jobrepo.StatusFailed
}
