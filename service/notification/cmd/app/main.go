package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmingruby/doris/lib/delivery/http/rest"
	"github.com/charmingruby/doris/lib/delivery/messaging/nats"
	"github.com/charmingruby/doris/lib/instrumentation/logger"
	"github.com/charmingruby/doris/service/notification/config"
	"github.com/charmingruby/doris/service/notification/internal/notification"
	"github.com/charmingruby/doris/service/notification/internal/platform"
	"github.com/gin-gonic/gin"
)

func main() {
	var log *logger.Logger

	cfg, err := config.New()
	if err != nil {
		log = logger.New(logger.LOG_LEVEL_INFO)
		log.Error("failed to load config", "error", err)
		return
	}

	log = logger.New(cfg.LogLevel)

	log.Info("config loaded successfully")

	sub, err := nats.NewSubscriber(
		log,
		nats.WithStream(cfg.Custom.NatsStream),
	)
	if err != nil {
		log.Error("failed to create nats subscriber", "error", err)
		return
	}

	log.Info("nats subscriber created successfully")

	server, router := rest.NewServer(cfg.Custom.RestServerHost, cfg.Custom.RestServerPort)

	initModules(log, cfg, router, sub)

	log.Info("modules initialized successfully")

	go func() {
		log.Info(fmt.Sprintf("rest server is running on %s:%s", cfg.Custom.RestServerHost, cfg.Custom.RestServerPort))

		if err := server.Start(); err != nil {
			log.Error("failed to start rest server", "error", err)
			return
		}
	}()

	gracefulShutdown(log, server, sub)
}

func initModules(log *logger.Logger, cfg config.Config, r *gin.Engine, sub *nats.Subscriber) {
	notification.NewEventHandler(log, sub, cfg)

	platform.NewHTTPHandler(r)
}

func gracefulShutdown(log *logger.Logger, srv *rest.Server, sub *nats.Subscriber) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sub.Close(ctx); err != nil {
		log.Error("failed to close nats subscriber", "error", err)
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown rest server", "error", err)
	}

	log.Info("shutdown complete")
}
