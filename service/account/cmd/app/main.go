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
	"github.com/charmingruby/doris/lib/persistence/postgres"
	"github.com/charmingruby/doris/lib/security"
	"github.com/charmingruby/doris/lib/validation"
	"github.com/charmingruby/doris/service/account/config"
	"github.com/charmingruby/doris/service/account/internal/access"
	"github.com/charmingruby/doris/service/account/internal/platform"

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

	db, err := postgres.New(logger, postgres.ConnectionInput{
		User:         cfg.Custom.DatabaseUser,
		Password:     cfg.Custom.DatabasePassword,
		Host:         cfg.Custom.DatabaseHost,
		DatabaseName: cfg.Custom.DatabaseName,
		SSL:          cfg.Custom.DatabaseSSL,
	})

	if err != nil {
		logger.Error("failed to create postgres connection", "error", err)
		return
	}

	logger.Info("postgres connection created successfully")

	server, router := rest.NewServer(cfg.Custom.RestServerHost, cfg.Custom.RestServerPort)

	val := validation.NewValidator()

	if err := initModules(logger, cfg, val, db, pub, router); err != nil {
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

	gracefulShutdown(logger, db, pub, server)
}

func initModules(logger *instrumentation.Logger, cfg config.Config, val *validation.Validator, db *postgres.Client, pub *nats.Publisher, r *gin.Engine) error {
	tokenClient := security.NewJWT(cfg.Custom.JWTIssuer, cfg.Custom.JWTSecret)

	mw := rest.NewMiddleware(tokenClient)

	accessDatasource, err := access.NewDatasource(db.Conn)
	if err != nil {
		return err
	}

	accessEvtHandler := access.NewEventHandler(pub, cfg)

	accessSvc := access.NewService(logger, accessDatasource, accessEvtHandler, tokenClient)

	access.NewHTTPHandler(logger, r, mw, val, accessSvc)

	platform.NewHTTPHandler(r)

	return nil
}

func gracefulShutdown(logger *instrumentation.Logger, db *postgres.Client, pub *nats.Publisher, srv *rest.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pub.Close(ctx); err != nil {
		logger.Error("failed to close nats publisher", "error", err)
	}

	if err := db.Close(ctx); err != nil {
		logger.Error("failed to disconnect from postgres", "error", err)
	}

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown rest server", "error", err)
	}

	logger.Info("shutdown complete")
}
