package usecase

import (
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/service/scribe/internal/codex/core/repository"
)

type UseCase struct {
	logger    *instrumentation.Logger
	codexRepo repository.CodexRepository
}

func New(logger *instrumentation.Logger, codexRepo repository.CodexRepository) *UseCase {
	return &UseCase{
		logger:    logger,
		codexRepo: codexRepo,
	}
}
