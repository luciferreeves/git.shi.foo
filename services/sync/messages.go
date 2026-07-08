package sync

const (
	EnqueueLog = "Failed to enqueue fetch job: %v"

	InvalidSignature    = "Invalid webhook signature."
	FetchMissingRepo    = "fetch job has no repo."
	MissingInstallation = "fetch job has no installation id."
)
