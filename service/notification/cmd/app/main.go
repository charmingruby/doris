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
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/service/notification/config"
	"github.com/charmingruby/doris/service/notification/internal/notification"
	"github.com/charmingruby/doris/service/notification/internal/platform"
	"github.com/charmingruby/doris/service/notification/test/memory"
	"github.com/gin-gonic/gin"
)

func main() {
	var logger *instrumentation.Logger

	cfg, err := config.New()
	if err != nil {
		logger = instrumentation.New(instrumentation.LOG_LEVEL_INFO)
		logger.Error("failed to load config", "error", err)
		return
	}

	logger = instrumentation.New(cfg.LogLevel)

	logger.Info("config loaded successfully")

	sub, err := nats.NewSubscriber(
		logger,
		nats.WithStream(cfg.Custom.NatsStream),
	)
	if err != nil {
		logger.Error("failed to create nats subscriber", "error", err)
		return
	}

	logger.Info("nats subscriber created successfully")

	server, router := rest.NewServer(cfg.Custom.RestServerHost, cfg.Custom.RestServerPort)

	initModules(logger, cfg, router, sub)

	logger.Info("modules initialized successfully")

	go func() {
		logger.Info(fmt.Sprintf("rest server is running on %s:%s", cfg.Custom.RestServerHost, cfg.Custom.RestServerPort))

		if err := server.Start(); err != nil {
			logger.Error("failed to start rest server", "error", err)
			return
		}
	}()

	gracefulShutdown(logger, server, sub)
}

func initModules(logger *instrumentation.Logger, cfg config.Config, r *gin.Engine, sub *nats.Subscriber) {
	notificationRepo := memory.NewNotificationRepository()

	notificationSvc := notification.NewService(logger, notificationRepo)

	notification.NewEventHandler(logger, sub, cfg, notificationSvc)

	platform.NewHTTPHandler(r)
}

func gracefulShutdown(logger *instrumentation.Logger, srv *rest.Server, sub *nats.Subscriber) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sub.Close(ctx); err != nil {
		logger.Error("failed to close nats subscriber", "error", err)
	}

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown rest server", "error", err)
	}

	logger.Info("shutdown complete")
}
