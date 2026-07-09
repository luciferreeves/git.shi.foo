package repos

import (
	"fmt"
	"net"
	"net/url"
	"strings"

	"git.shi.foo/config"
)

func sshCloneURL(owner string, name string) string {
	host, port := sshHostPort()

	if port == SSHDefaultPort {
		return fmt.Sprintf(SSHCloneScpFormat, SSHUser, host, owner, name)
	}
	return fmt.Sprintf(SSHCloneURLFormat, SSHUser, host, port, owner, name)
}

func sshHostPort() (string, string) {
	host := SSHFallbackHost
	if parsed, parseError := url.Parse(config.Server.PublicURL); parseError == nil && parsed.Hostname() != "" {
		host = parsed.Hostname()
	}

	port := SSHDefaultPort
	if _, listenPort, splitError := net.SplitHostPort(config.Git.SSHAddress); splitError == nil && listenPort != "" {
		port = listenPort
	} else if trimmed := strings.TrimPrefix(config.Git.SSHAddress, ":"); trimmed != "" {
		port = trimmed
	}

	return host, port
}
