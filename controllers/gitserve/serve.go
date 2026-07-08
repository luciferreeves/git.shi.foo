package gitserve

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"git.shi.foo/git"
	servicegitserve "git.shi.foo/services/gitserve"
	"git.shi.foo/utils/logger"
	"git.shi.foo/utils/shortcuts"

	"github.com/gofiber/fiber/v2"
)

func InfoRefs(context *fiber.Ctx) error {
	service := context.Query("service")
	owner := context.Params("owner")
	name := repoName(context)

	currentUser := servicegitserve.ResolveUser(basicPassword(context))
	if authError := servicegitserve.Authorize(context.UserContext(), currentUser, owner, name, service); authError != nil {
		return prompt(context, authError)
	}

	advertisement, adviseError := servicegitserve.Advertise(service, owner, name)
	if adviseError != nil {
		return adviseError
	}

	context.Set(fiber.HeaderContentType, fmt.Sprintf(AdvertisementType, service))
	context.Set(fiber.HeaderCacheControl, NoCache)
	return context.Send(advertisement)
}

func UploadPack(context *fiber.Ctx) error {
	return servePack(context, git.ServiceUploadPack)
}

func ReceivePack(context *fiber.Ctx) error {
	return servePack(context, git.ServiceReceivePack)
}

func servePack(context *fiber.Ctx, service string) error {
	owner := context.Params("owner")
	name := repoName(context)

	currentUser := servicegitserve.ResolveUser(basicPassword(context))
	if authError := servicegitserve.Authorize(context.UserContext(), currentUser, owner, name, service); authError != nil {
		return prompt(context, authError)
	}

	body, bodyError := requestBody(context)
	if bodyError != nil {
		return shortcuts.ServiceError(fiber.StatusBadRequest, InvalidBody)
	}

	accessToken := ""
	if service == git.ServiceReceivePack {
		token, tokenError := servicegitserve.TokenFor(context.UserContext(), currentUser.ID)
		if tokenError != nil {
			return tokenError
		}
		accessToken = token
	}

	context.Set(fiber.HeaderContentType, fmt.Sprintf(ResultType, service))
	context.Set(fiber.HeaderCacheControl, NoCache)

	context.Context().SetBodyStreamWriter(func(writer *bufio.Writer) {
		var runError error
		if service == git.ServiceUploadPack {
			runError = git.UploadPack(owner, name, body, writer)
		} else {
			runError = git.ReceivePack(owner, name, accessToken, body, writer)
		}
		if runError != nil {
			logger.Errorf(LogPrefix, ServiceLog, runError)
		}
		writer.Flush()
	})

	return nil
}

func repoName(context *fiber.Ctx) string {
	return strings.TrimSuffix(context.Params("repo"), git.RepoSuffix)
}

func requestBody(context *fiber.Ctx) (io.Reader, error) {
	body := context.Body()
	if context.Get(fiber.HeaderContentEncoding) == GzipEncoding {
		return gzip.NewReader(bytes.NewReader(body))
	}
	return bytes.NewReader(body), nil
}

func basicPassword(context *fiber.Ctx) string {
	header := context.Get(fiber.HeaderAuthorization)
	if !strings.HasPrefix(header, BasicPrefix) {
		return ""
	}

	decoded, decodeError := base64.StdEncoding.DecodeString(strings.TrimPrefix(header, BasicPrefix))
	if decodeError != nil {
		return ""
	}

	parts := strings.SplitN(string(decoded), CredentialColon, 2)
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

func prompt(context *fiber.Ctx, authError *fiber.Error) error {
	if authError.Code == fiber.StatusUnauthorized {
		context.Set(fiber.HeaderWWWAuthenticate, AuthRealm)
	}
	return authError
}
