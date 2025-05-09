package usecase

import (
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/lib/persistence"
	"github.com/charmingruby/doris/lib/security"
	"github.com/charmingruby/doris/service/account/internal/access/core/event"
	"github.com/charmingruby/doris/service/account/internal/access/core/repository"
)

type UseCase struct {
	logger      *instrumentation.Logger
	apiKeyRepo  repository.APIKeyRepository
	otpRepo     repository.OTPRepository
	txManager   persistence.TransactionManager[repository.TransactionManager]
	tokenClient security.Token
	event       event.Handler
}

func New(
	logger *instrumentation.Logger,
	apiKeyRepo repository.APIKeyRepository,
	otpRepo repository.OTPRepository,
	txManager persistence.TransactionManager[repository.TransactionManager],
	tokenClient security.Token,
	event event.Handler,
) *UseCase {
	return &UseCase{
		logger:      logger,
		apiKeyRepo:  apiKeyRepo,
		otpRepo:     otpRepo,
		txManager:   txManager,
		tokenClient: tokenClient,
		event:       event,
	}
}
