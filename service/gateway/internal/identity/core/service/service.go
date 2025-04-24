package service

import (
	"github.com/charmingruby/doris/lib/instrumentation/logger"
	"github.com/charmingruby/doris/service/gateway/internal/identity/core/repository"
	"github.com/charmingruby/doris/service/gateway/internal/identity/delivery/event"
)

type Service struct {
	logger       *logger.Logger
	apiKeyRepo   repository.APIKeyRepository
	eventHandler *event.Handler
}

func New(
	logger *logger.Logger,
	apiKeyRepo repository.APIKeyRepository,
	eventHandler *event.Handler,
) *Service {
	return &Service{
		logger:       logger,
		apiKeyRepo:   apiKeyRepo,
		eventHandler: eventHandler,
	}
}
