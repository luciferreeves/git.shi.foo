package sync

const (
	EnqueueLog                 = "Failed to enqueue fetch job: %v."
	ReconcileLog               = "Failed to reconcile repository mirrors: %v."
	RepositoryIdentityMismatch = "GitHub returned a different repository for %s/%s. Synchronization was refused."

	InvalidSignature    = "Invalid webhook signature."
	FetchMissingRepo    = "fetch job has no repo."
	MissingInstallation = "fetch job has no installation id."
)
