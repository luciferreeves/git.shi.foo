package ssh

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"os"

	"git.shi.foo/config"
	"git.shi.foo/utils/logger"

	gossh "golang.org/x/crypto/ssh"
)

func hostSigner() (gossh.Signer, error) {
	path := config.Git.SSHHostKey
	if path != "" {
		if data, readError := os.ReadFile(path); readError == nil {
			return gossh.ParsePrivateKey(data)
		}
	}

	_, privateKey, generateError := ed25519.GenerateKey(rand.Reader)
	if generateError != nil {
		return nil, generateError
	}

	signer, signerError := gossh.NewSignerFromKey(privateKey)
	if signerError != nil {
		return nil, signerError
	}

	if path == "" {
		logger.Infof(LogPrefix, HostKeyGenerated)
		return signer, nil
	}

	persistHostKey(path, privateKey)
	return signer, nil
}

func persistHostKey(path string, privateKey ed25519.PrivateKey) {
	block, marshalError := gossh.MarshalPrivateKey(privateKey, "")
	if marshalError != nil {
		logger.Infof(LogPrefix, HostKeyGenerated)
		return
	}

	if writeError := os.WriteFile(path, pem.EncodeToMemory(block), HostKeyPermission); writeError != nil {
		logger.Infof(LogPrefix, HostKeyGenerated)
		return
	}

	logger.Infof(LogPrefix, HostKeyPersisted, path)
}
