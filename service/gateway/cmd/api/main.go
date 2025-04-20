package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmingruby/doris/lib/delivery/rest"
	"github.com/charmingruby/doris/lib/instrumentation/logger"
	"github.com/charmingruby/doris/service/gateway/config"
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

	server, _ := rest.NewServer(cfg.Custom.RestServerHost, cfg.Custom.RestServerPort)

	go func() {
		if err := server.Start(); err != nil {
			log.Error("failed to start rest server", "error", err)
			return
		}
	}()

	gracefulShutdown(log, server)
}

func gracefulShutdown(log *logger.Logger, srv *rest.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown rest server", "error", err)
	}
}
