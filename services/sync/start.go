package sync

import (
	"context"
	"time"

	"git.shi.foo/utils/logger"
)

func Start(runContext context.Context) {
	go func() {
		ticker := time.NewTicker(ReconcileInterval)
		defer ticker.Stop()

		for {
			if runContext.Err() != nil {
				return
			}
			if reconcileError := Reconcile(); reconcileError != nil {
				logger.Errorf(LogPrefix, ReconcileLog, reconcileError)
			}
			select {
			case <-runContext.Done():
				return
			case <-ticker.C:
			}
		}
	}()
}
