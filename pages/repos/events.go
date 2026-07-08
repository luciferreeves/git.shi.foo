package repos

import (
	"bufio"
	"encoding/json"
	"time"

	"git.shi.foo/events"
	"git.shi.foo/jobs"
	"git.shi.foo/services/repos"
	"git.shi.foo/utils/meta"

	"github.com/gofiber/fiber/v2"
)

func ImportEvents(context *fiber.Ctx) error {
	streamContext, streamError := repos.ImportStream(context.UserContext(), meta.User(context), context.Params("owner"), context.Params("repo"))
	if streamError != nil {
		return streamError
	}

	context.Set(fiber.HeaderContentType, SSEContentType)
	context.Set(fiber.HeaderCacheControl, SSECacheControl)
	context.Set(fiber.HeaderConnection, SSEConnection)

	repoID := streamContext.RepoID
	topic := streamContext.Topic

	context.Context().SetBodyStreamWriter(func(writer *bufio.Writer) {
		var channel <-chan events.Event
		if topic != "" {
			subscription, unsubscribe := events.Default.Subscribe(topic)
			defer unsubscribe()
			channel = subscription
		}

		snapshot := repos.ImportSnapshot(repoID)
		if writeError := writeProgress(writer, snapshot); writeError != nil {
			return
		}
		if snapshot.Done || channel == nil {
			return
		}

		heartbeat := time.NewTicker(SSEHeartbeatInterval)
		defer heartbeat.Stop()

		for {
			select {
			case event, open := <-channel:
				if !open {
					return
				}
				progress, isProgress := event.Data.(jobs.Progress)
				if !isProgress {
					continue
				}
				if writeError := writeProgress(writer, progress); writeError != nil {
					return
				}
				if progress.Done {
					return
				}
			case <-heartbeat.C:
				if _, writeError := writer.WriteString(SSEHeartbeat); writeError != nil {
					return
				}
				if flushError := writer.Flush(); flushError != nil {
					return
				}
			}
		}
	})

	return nil
}

func writeProgress(writer *bufio.Writer, progress jobs.Progress) error {
	payload, marshalError := json.Marshal(progress)
	if marshalError != nil {
		return marshalError
	}

	if _, writeError := writer.WriteString(SSEDataPrefix + string(payload) + SSEFrameSuffix); writeError != nil {
		return writeError
	}

	return writer.Flush()
}
