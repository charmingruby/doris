package service

import (
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/service/hub/internal/identity/core/event"
	"github.com/charmingruby/doris/service/hub/internal/identity/core/repository"
)

type Service struct {
	logger *instrumentation.Logger
	repo   repository.APIKeyRepository
	event  event.Handler
}

func New(
	logger *instrumentation.Logger,
	repo repository.APIKeyRepository,
	event event.Handler,
) *Service {
	return &Service{
		logger: logger,
		repo:   repo,
		event:  event,
	}
}
