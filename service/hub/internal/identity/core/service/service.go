package service

import (
	"github.com/charmingruby/doris/lib/instrumentation/logger"
	"github.com/charmingruby/doris/service/hub/internal/identity/core/event"
	"github.com/charmingruby/doris/service/hub/internal/identity/core/repository"
)

type Service struct {
	logger       *logger.Logger
	apiKeyRepo   repository.APIKeyRepository
	eventHandler event.Handler
}

func New(
	logger *logger.Logger,
	apiKeyRepo repository.APIKeyRepository,
	eventHandler event.Handler,
) *Service {
	return &Service{
		logger:       logger,
		apiKeyRepo:   apiKeyRepo,
		eventHandler: eventHandler,
	}
}
