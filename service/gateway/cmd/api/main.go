package main

import (
	"fmt"

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

	fmt.Printf("config: %+v\n", cfg)
}
