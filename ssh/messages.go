package ssh

const (
	Listening        = "SSH server listening on %s."
	ServeFailed      = "SSH serve failed: %v"
	ServiceLog       = "SSH git service failed: %v"
	HostKeyGenerated = "Generated ephemeral SSH host key (set SSH_HOST_KEY to persist)."
	HostKeyPersisted = "Persisted generated SSH host key to %s."

	InvalidCommand = "unsupported command\n"
	InvalidRepo    = "invalid repository\n"
	UnknownService = "unknown service\n"
	LineBreak      = "\n"
)
