package repos

import "time"

const (
	LogPrefix = "Repos"

	SSEContentType  = "text/event-stream"
	SSECacheControl = "no-cache"
	SSEConnection   = "keep-alive"

	SSEDataPrefix        = "data: "
	SSEFrameSuffix       = "\n\n"
	SSEHeartbeat         = ": ping\n\n"
	SSEHeartbeatInterval = 20 * time.Second
)
