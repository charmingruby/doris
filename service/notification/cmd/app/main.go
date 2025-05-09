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
	"github.com/charmingruby/doris/lib/persistence/dynamo"
	"github.com/charmingruby/doris/lib/security"
	"github.com/charmingruby/doris/lib/validation"
	"github.com/charmingruby/doris/service/notification/config"
	"github.com/charmingruby/doris/service/notification/internal/notification"
	"github.com/charmingruby/doris/service/notification/internal/notification/provider/notifier"
	"github.com/charmingruby/doris/service/notification/internal/platform"
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

	db, err := dynamo.New(logger, dynamo.ConnectionInput{
		Region: cfg.Custom.AWSRegion,
	})

	if err != nil {
		logger.Error("failed to create dynamo connection", "error", err)
		return
	}

	logger.Info("dynamo connection created successfully")

	val := validation.NewValidator()

	server, router := rest.NewServer(cfg.Custom.RestServerHost, cfg.Custom.RestServerPort)

	if err := initModules(logger, cfg, db, router, sub, val); err != nil {
		logger.Error("failed to initialize modules", "error", err)
		return
	}

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

func initModules(logger *instrumentation.Logger, cfg config.Config, db *dynamo.Client, r *gin.Engine, sub *nats.Subscriber, val *validation.Validator) error {
	notificationDatasource, err := notification.NewDatasource(cfg, db)
	if err != nil {
		return err
	}

	notifier, err := notifier.NewSES(cfg.Custom.AWSRegion, cfg.Custom.SourceEmail)
	if err != nil {
		return err
	}

	tokenClient := security.NewJWT(cfg.Custom.JWTIssuer, cfg.Custom.JWTSecret)

	mw := rest.NewMiddleware(tokenClient)

	notificationUc := notification.NewUseCase(logger, notificationDatasource, notifier)

	notification.NewEventHandler(logger, sub, cfg, notificationUc)

	notification.NewHTTPHandler(logger, r, mw, val, notificationUc)

	platform.NewHTTPHandler(r)

	return nil
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
