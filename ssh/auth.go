package ssh

import (
	keyrepo "git.shi.foo/repositories/key"

	glssh "github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

type contextKey string

const contextUserID contextKey = "user_id"

func handlePublicKey(sshContext glssh.Context, offeredKey glssh.PublicKey) bool {
	fingerprint := gossh.FingerprintSHA256(offeredKey)

	record, findError := keyrepo.FindByFingerprint(fingerprint)
	if findError != nil {
		return false
	}

	sshContext.SetValue(contextUserID, record.UserID)
	return true
}
