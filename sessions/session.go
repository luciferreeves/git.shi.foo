package sessions

import (
	"encoding/gob"
	"fmt"

	"git.shi.foo/config"
	"git.shi.foo/utils/collections"
	"git.shi.foo/utils/logger"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres/v3"
)

var Store *session.Store

func init() {
	gob.Register(collections.Record[string, any]{})
	logger.Infof(LogPrefix, GobTypeRecordRegistered)

	storage := postgres.New(postgres.Config{
		ConnectionURI: config.Database.DSN,
		Table:         config.Session.CookieName,
		Reset:         false,
		GCInterval:    SessionInterval,
	})

	Store = session.New(session.Config{
		Storage:        storage,
		Expiration:     config.Session.CookieTimeout,
		KeyLookup:      fmt.Sprintf("cookie:%s", config.Session.CookieName),
		CookieDomain:   config.Session.CookieDomain,
		CookiePath:     config.Session.CookiePath,
		CookieSecure:   config.Session.CookieSecure,
		CookieSameSite: config.Session.CookieSameSite,
		CookieHTTPOnly: true,
	})

	logger.Successf(LogPrefix, SessionStoreInitialized)
}
