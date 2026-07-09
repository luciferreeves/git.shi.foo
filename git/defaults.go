package git

import "os"

const (
	LogPrefix           = "Git"
	DirectoryPermission = os.FileMode(0o755)
	RepoSuffix          = ".git"
	CleanCloneURLFormat = "https://github.com/%s/%s.git"
	HeadRef             = "HEAD"
	CommitFormat        = "--format=%H%x1f%h%x1f%s%x1f%an%x1f%aI"
	TypeTree            = "tree"
	TypeBlob            = "blob"
	GitTokenEnv         = "GIT_ASKPASS_TOKEN"
	CredentialHelper    = "credential.helper=!f() { echo username=x-access-token; echo password=$GIT_ASKPASS_TOKEN; }; f"
	RedactedToken       = "[redacted]"

	PhaseCounting    = "counting"
	PhaseCompressing = "compressing"
	PhaseReceiving   = "receiving"
	PhaseResolving   = "resolving"

	ServiceUploadPack  = "git-upload-pack"
	ServiceReceivePack = "git-receive-pack"
	ServicePrefix      = "# service="
	FlushPacket        = "0000"

	HooksSubdir        = ".githooks"
	PreReceiveHookName = "pre-receive"
	HookPermission     = os.FileMode(0o755)
)
