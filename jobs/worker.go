package jobs

import (
	"context"
	"time"

	jobrepo "git.shi.foo/repositories/job"
	"git.shi.foo/utils/logger"
)

var wakeup = make(chan struct{}, 1)

func Signal() {
	select {
	case wakeup <- struct{}{}:
	default:
	}
}

func Start(runContext context.Context) {
	for range WorkerCount {
		go workerLoop(runContext)
	}
}

func workerLoop(runContext context.Context) {
	ticker := time.NewTicker(PollInterval)
	defer ticker.Stop()

	for {
		record, claimError := jobrepo.ClaimNext(time.Now())
		if claimError != nil {
			logger.Errorf(LogPrefix, ClaimFailed, claimError)
		} else if record != nil {
			Signal()
			process(runContext, record)
			continue
		}

		select {
		case <-runContext.Done():
			return
		case <-wakeup:
		case <-ticker.C:
		}
	}
}
