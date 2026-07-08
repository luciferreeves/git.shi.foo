package jobs

import "time"

const (
	LogPrefix = "Jobs"

	WorkerCount      = 3
	PollInterval     = 15 * time.Second
	RetryBackoffBase = 10 * time.Second
	PersistThreshold = 10

	TopicFormat   = "job:%d"
	EventProgress = "progress"
)
