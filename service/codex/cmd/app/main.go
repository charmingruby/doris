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
	"github.com/charmingruby/doris/lib/storage/s3"
	"github.com/charmingruby/doris/lib/validation"
	"github.com/charmingruby/doris/service/codex/config"
	"github.com/charmingruby/doris/service/codex/internal/codex"
	"github.com/charmingruby/doris/service/codex/internal/platform"
	"github.com/charmingruby/doris/service/codex/internal/quota"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
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
		Port:         cfg.Custom.DatabasePort,
		DatabaseName: cfg.Custom.DatabaseName,
		SSL:          cfg.Custom.DatabaseSSL,
	})
	if err != nil {
		logger.Error("failed to create postgres connection", "error", err)
		return
	}

	logger.Info("postgres connection created successfully")

	val := validation.NewValidator()

	server, router := rest.NewServer(cfg.Custom.RestServerHost, cfg.Custom.RestServerPort)

	resetAllQuotaUsagesFn, err := initModules(logger, cfg, router, db, sub, pub, val)
	if err != nil {
		logger.Error("failed to initialize modules", "error", err)
		return
	}

	monthlyQuotaUsageResetCronJob(logger, resetAllQuotaUsagesFn)

	logger.Info("modules initialized successfully")

	go func() {
		logger.Info(fmt.Sprintf("rest server is running on %s:%s", cfg.Custom.RestServerHost, cfg.Custom.RestServerPort))

		if err := server.Start(); err != nil {
			logger.Error("failed to start rest server", "error", err)
			return
		}
	}()

	gracefulShutdown(logger, server, sub, db)
}

func initModules(
	logger *instrumentation.Logger,
	cfg config.Config,
	r *gin.Engine,
	db *postgres.Client,
	sub *nats.Subscriber,
	pub *nats.Publisher,
	val *validation.Validator,
) (func(context.Context) error, error) {
	// Shared
	tokenClient := security.NewJWT(cfg.Custom.JWTIssuer, cfg.Custom.JWTSecret)
	mw := rest.NewMiddleware(tokenClient)

	storage, err := s3.New(logger, cfg.Custom.AWSRegion)
	if err != nil {
		return nil, err
	}

	platform.NewHTTPHandler(r)

	// Quota
	quotaDatasource, err := quota.NewDatasource(db.Conn)
	if err != nil {
		return nil, err
	}

	quotaUseCase := quota.NewUseCase(logger, quotaDatasource)

	quotaProvider := quota.NewProvider(logger, quotaDatasource)

	quota.NewHTTPHandler(logger, r, mw, val, quotaUseCase)

	quota.NewEventHandler(logger, sub, cfg, quotaUseCase)

	// Codex
	codexDatasource, err := codex.NewDatasource(db.Conn)
	if err != nil {
		return nil, err
	}

	codexEventHandler, err := codex.NewEventHandler(logger, pub, sub, cfg)
	if err != nil {
		return nil, err
	}

	codexUseCase := codex.NewUseCase(logger, codexDatasource, codexEventHandler, storage, quotaProvider.QuotaUsageManagement, cfg.Custom.AWSEmbeddingSourceDocsBucket)

	codex.NewHTTPHandler(logger, r, mw, val, codexUseCase)

	return quotaUseCase.ResetAllQuotaUsages, nil
}

func gracefulShutdown(logger *instrumentation.Logger, srv *rest.Server, sub *nats.Subscriber, db *postgres.Client) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sub.Close(ctx); err != nil {
		logger.Error("failed to close nats subscriber", "error", err)
	}

	if err := db.Close(ctx); err != nil {
		logger.Error("failed to disconnect from postgres", "error", err)
	}

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown rest server", "error", err)
	}

	logger.Info("shutdown complete")
}

func monthlyQuotaUsageResetCronJob(logger *instrumentation.Logger, fn func(context.Context) error) {
	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc("0 0 0 1 * *", func() {
		logger.Info("cron job: reset all the quota usages", "time", time.Now())

		if err := fn(context.Background()); err != nil {
			logger.Error(err.Error())
		} else {
			logger.Info("cron job: reset all the quota changes successfully")
		}
	})

	if err != nil {
		logger.Error("failed to register cron job", "err", err)
		return
	}

	c.Start()
}
