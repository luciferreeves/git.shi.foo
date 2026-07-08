package jobs

const (
	ClaimFailed   = "Failed to claim job: %v"
	RequeueFailed = "Failed to requeue running jobs: %v"
	PersistFailed = "Failed to persist job: %v"
	PanicLog      = "Job %d panicked: %v"

	UnknownKind = "no runner registered for job kind %s"
	JobPanicked = "job panicked: %v"
)
