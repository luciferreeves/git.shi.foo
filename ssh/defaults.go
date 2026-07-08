package ssh

import "os"

const (
	LogPrefix = "SSH"

	HostKeyPermission = os.FileMode(0o600)
	HostKeyDirMode    = os.FileMode(0o700)
	PathSeparator     = "/"
)
