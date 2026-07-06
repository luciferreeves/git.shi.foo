package sessions

import "time"

const (
	LogPrefix       = "Session"
	SessionInterval = 10 * time.Second
	SessionAuthKey  = "provider_id"
	SessionLocalKey = "Session"
)
