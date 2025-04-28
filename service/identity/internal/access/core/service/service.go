package service

import (
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/service/identity/internal/access/core/event"
	"github.com/charmingruby/doris/service/identity/internal/access/core/repository"
)

type Service struct {
	logger     *instrumentation.Logger
	apiKeyRepo repository.APIKeyRepository
	otpRepo    repository.OTPRepository
	event      event.Handler
}

func New(
	logger *instrumentation.Logger,
	apiKeyRepo repository.APIKeyRepository,
	otpRepo repository.OTPRepository,
	event event.Handler,
) *Service {
	return &Service{
		logger:     logger,
		apiKeyRepo: apiKeyRepo,
		otpRepo:    otpRepo,
		event:      event,
	}
}
