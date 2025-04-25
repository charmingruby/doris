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
	"github.com/charmingruby/doris/lib/persistence/mongo"
	"github.com/charmingruby/doris/lib/validation"
	"github.com/charmingruby/doris/service/gateway/config"
	"github.com/charmingruby/doris/service/gateway/internal/identity"
	"github.com/charmingruby/doris/service/gateway/internal/platform"
	"github.com/charmingruby/doris/service/gateway/test/memory"
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

	pub, err := nats.NewPublisher(
		log,
		nats.WithStream(cfg.Custom.NatsStream),
	)
	if err != nil {
		log.Error("failed to create nats publisher", "error", err)
		return
	}

	log.Info("nats publisher created successfully")

	db, err := mongo.New(cfg.Custom.MongoURL, cfg.Custom.MongoDatabase)
	if err != nil {
		log.Error("failed to create mongo connection", "error", err)
		return
	}

	log.Info("mongo connection created successfully")

	server, router := rest.NewServer(cfg.Custom.RestServerHost, cfg.Custom.RestServerPort)

	val := validation.NewValidator()

	initModules(log, cfg, val, db, pub, router)

	log.Info("modules initialized successfully")

	go func() {
		log.Info(fmt.Sprintf("rest server is running on %s:%s", cfg.Custom.RestServerHost, cfg.Custom.RestServerPort))

		if err := server.Start(); err != nil {
			log.Error("failed to start rest server", "error", err)
			return
		}
	}()

	gracefulShutdown(log, db, pub, server)
}

func initModules(log *logger.Logger, cfg config.Config, val *validation.Validator, db *mongo.Client, pub *nats.Publisher, r *gin.Engine) {
	apiKeyRepo := memory.NewAPIKeyRepository()

	identityEvtHandler := identity.NewEventHandler(pub, cfg)

	identitySvc := identity.NewService(log, apiKeyRepo, identityEvtHandler)

	identity.NewHTTPHandler(log, r, val, identitySvc)

	platform.NewHTTPHandler(r)
}

func gracefulShutdown(log *logger.Logger, db *mongo.Client, pub *nats.Publisher, srv *rest.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pub.Close(ctx); err != nil {
		log.Error("failed to close nats publisher", "error", err)
	}

	if err := db.Close(ctx); err != nil {
		log.Error("failed to disconnect from mongo", "error", err)
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown rest server", "error", err)
	}

	log.Info("shutdown complete")
}
