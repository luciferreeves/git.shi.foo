package jobs

import (
	jobrepo "git.shi.foo/repositories/job"
	"git.shi.foo/utils/logger"
)

func Recover() {
	if requeueError := jobrepo.RequeueRunning(); requeueError != nil {
		logger.Errorf(LogPrefix, RequeueFailed, requeueError)
	}
}
