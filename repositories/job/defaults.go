package job

const (
	LogPrefix = "JobRepository"

	KindImport = "import"
	KindFetch  = "fetch"

	StatusPending   = "pending"
	StatusRunning   = "running"
	StatusSucceeded = "succeeded"
	StatusFailed    = "failed"
)
