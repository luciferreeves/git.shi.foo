package sync

import "time"

const (
	LogPrefix = "Sync"

	EventPush   = "push"
	EventCreate = "create"
	EventDelete = "delete"

	FetchMaxAttempts  = 3
	ReconcileInterval = 5 * time.Minute
)
