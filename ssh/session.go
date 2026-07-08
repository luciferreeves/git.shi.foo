package ssh

import (
	"io"
	"strings"

	"git.shi.foo/account"
	"git.shi.foo/git"
	userrepo "git.shi.foo/repositories/user"
	"git.shi.foo/services/gitserve"
	"git.shi.foo/utils/logger"

	glssh "github.com/gliderlabs/ssh"
)

func handleSession(session glssh.Session) {
	command := session.Command()
	if len(command) < 2 {
		io.WriteString(session.Stderr(), InvalidCommand)
		session.Exit(1)
		return
	}

	service := command[0]
	owner, name, valid := parseRepoArg(command[1])
	if !valid {
		io.WriteString(session.Stderr(), InvalidRepo)
		session.Exit(1)
		return
	}

	userID, _ := session.Context().Value(contextUserID).(uint)
	currentUser := loadUser(userID)

	if authError := gitserve.Authorize(session.Context(), currentUser, owner, name, service); authError != nil {
		io.WriteString(session.Stderr(), authError.Message+LineBreak)
		session.Exit(1)
		return
	}

	switch service {
	case git.ServiceUploadPack:
		if runError := git.UploadPackSession(owner, name, session, session); runError != nil {
			logger.Errorf(LogPrefix, ServiceLog, runError)
			session.Exit(1)
			return
		}
	case git.ServiceReceivePack:
		accessToken, tokenError := gitserve.TokenFor(session.Context(), userID)
		if tokenError != nil {
			io.WriteString(session.Stderr(), tokenError.Message+LineBreak)
			session.Exit(1)
			return
		}
		if runError := git.ReceivePackSession(owner, name, accessToken, session, session); runError != nil {
			logger.Errorf(LogPrefix, ServiceLog, runError)
			session.Exit(1)
			return
		}
	default:
		io.WriteString(session.Stderr(), UnknownService)
		session.Exit(1)
		return
	}

	session.Exit(0)
}

func parseRepoArg(argument string) (string, string, bool) {
	trimmed := strings.TrimSuffix(strings.TrimPrefix(argument, PathSeparator), git.RepoSuffix)
	parts := strings.Split(trimmed, PathSeparator)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func loadUser(userID uint) *account.Response {
	record, findError := userrepo.FindByID(userID)
	if findError != nil {
		return nil
	}

	response := record.ToResponse()
	return &response
}
