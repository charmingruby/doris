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
	"github.com/charmingruby/doris/lib/persistence/mongo"
	"github.com/charmingruby/doris/lib/validation"
	"github.com/charmingruby/doris/service/hub/config"
	"github.com/charmingruby/doris/service/hub/internal/identity"
	"github.com/charmingruby/doris/service/hub/internal/platform"
	"github.com/charmingruby/doris/service/hub/test/memory"

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

	pub, err := nats.NewPublisher(
		logger,
		nats.WithStream(cfg.Custom.NatsStream),
	)
	if err != nil {
		logger.Error("failed to create nats publisher", "error", err)
		return
	}

	logger.Info("nats publisher created successfully")

	db, err := mongo.New(cfg.Custom.MongoURL, cfg.Custom.MongoDatabase)
	if err != nil {
		logger.Error("failed to create mongo connection", "error", err)
		return
	}

	logger.Info("mongo connection created successfully")

	server, router := rest.NewServer(cfg.Custom.RestServerHost, cfg.Custom.RestServerPort)

	val := validation.NewValidator()

	initModules(logger, cfg, val, db, pub, router)

	logger.Info("modules initialized successfully")

	go func() {
		logger.Info(fmt.Sprintf("rest server is running on %s:%s", cfg.Custom.RestServerHost, cfg.Custom.RestServerPort))

		if err := server.Start(); err != nil {
			logger.Error("failed to start rest server", "error", err)
			return
		}
	}()

	gracefulShutdown(logger, db, pub, server)
}

func initModules(logger *instrumentation.Logger, cfg config.Config, val *validation.Validator, db *mongo.Client, pub *nats.Publisher, r *gin.Engine) {
	apiKeyRepo := memory.NewAPIKeyRepository()

	otpRepo := memory.NewOTPRepository()

	identityEvtHandler := identity.NewEventHandler(pub, cfg)

	identitySvc := identity.NewService(logger, apiKeyRepo, otpRepo, identityEvtHandler)

	identity.NewHTTPHandler(logger, r, val, identitySvc)

	platform.NewHTTPHandler(r)
}

func gracefulShutdown(logger *instrumentation.Logger, db *mongo.Client, pub *nats.Publisher, srv *rest.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pub.Close(ctx); err != nil {
		logger.Error("failed to close nats publisher", "error", err)
	}

	if err := db.Close(ctx); err != nil {
		logger.Error("failed to disconnect from mongo", "error", err)
	}

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown rest server", "error", err)
	}

	logger.Info("shutdown complete")
}
