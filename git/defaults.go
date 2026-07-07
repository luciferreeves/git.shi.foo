package git

import "os"

const (
	LogPrefix           = "Git"
	DirectoryPermission = os.FileMode(0o755)
	RepoSuffix          = ".git"
	CleanCloneURLFormat = "https://github.com/%s/%s.git"
	GitTokenEnv         = "GIT_ASKPASS_TOKEN"
	CredentialHelper    = "credential.helper=!f() { echo username=x-access-token; echo password=$GIT_ASKPASS_TOKEN; }; f"
	RedactedToken       = "[redacted]"
)
