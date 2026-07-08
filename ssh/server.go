package ssh

import (
	"context"

	"git.shi.foo/config"
	"git.shi.foo/utils/logger"

	glssh "github.com/gliderlabs/ssh"
)

func Start(runContext context.Context) error {
	signer, signerError := hostSigner()
	if signerError != nil {
		return signerError
	}

	server := &glssh.Server{
		Addr:             config.Git.SSHAddress,
		Handler:          handleSession,
		PublicKeyHandler: handlePublicKey,
	}
	server.AddHostKey(signer)

	go func() {
		<-runContext.Done()
		server.Close()
	}()

	go func() {
		logger.Successf(LogPrefix, Listening, config.Git.SSHAddress)
		if serveError := server.ListenAndServe(); serveError != nil && serveError != glssh.ErrServerClosed {
			logger.Errorf(LogPrefix, ServeFailed, serveError)
		}
	}()

	return nil
}
