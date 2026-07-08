package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"git.shi.foo/config"
	"git.shi.foo/jobs"
	"git.shi.foo/middleware"
	jobrepo "git.shi.foo/repositories/job"
	"git.shi.foo/router"
	"git.shi.foo/services/repos"
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

	jobs.Register(jobrepo.KindImport, repos.RunImport)
	jobs.Recover()
	repos.ReconcileImports()

	workerContext, stopWorkers := context.WithCancel(context.Background())
	jobs.Start(workerContext)

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
	stopWorkers()

	if shutdownError := application.Shutdown(); shutdownError != nil {
		logger.Errorf(LogPrefix, ServerShutdownFailed, shutdownError)
	}

	logger.Successf(LogPrefix, ServerShutdownComplete)
}
