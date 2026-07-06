package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"git.shi.foo/config"
	"git.shi.foo/middleware"
	"git.shi.foo/router"
	"git.shi.foo/tags"
	"git.shi.foo/utils/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/django/v3"
)

func main() {
	tags.Initialize()
	engine := django.New("./templates", ".django")
	engine.Reload(config.Server.Debug)

	application := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		BodyLimit:             32 * 1024 * 1024,
		Views:                 engine,
		ErrorHandler:          router.ErrorHandler,
	})

	application.Use(recover.New())

	middleware.Initialize(application)
	router.Initialize(application)

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		address := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
		logger.Successf(LogPrefix, ServerStarting, address)

		if listenError := application.Listen(address); listenError != nil {
			logger.Fatalf(LogPrefix, ServerListenFailed, listenError)
		}
	}()

	<-shutdownSignal
	logger.Infof(LogPrefix, ServerShuttingDown)

	if shutdownError := application.Shutdown(); shutdownError != nil {
		logger.Errorf(LogPrefix, ServerShutdownFailed, shutdownError)
	}

	logger.Successf(LogPrefix, ServerShutdownComplete)
}
