package git

const (
	MkdirFailed       = "failed to create repo directory: %v"
	CloneFailed       = "git clone failed: %v: %s"
	GitCommandFailed  = "git command failed: %v: %s"
	CommitParseFailed = "could not parse commit output"
	UnknownService    = "unknown git service: %s"
	ServiceFailed     = "git service failed: %v: %s"
	HooksDirFailed    = "failed to create hooks directory: %v"
	HookWriteFailed   = "failed to write pre-receive hook: %v"
)
