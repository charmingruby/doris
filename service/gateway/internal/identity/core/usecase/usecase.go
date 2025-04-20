package usecase

import (
	"github.com/charmingruby/doris/service/gateway/internal/identity/core/repository"
)

type UseCase struct {
	apiKeyRepo repository.APIKeyRepository
}

func New(
	apiKeyRepo repository.APIKeyRepository,

) *UseCase {
	return &UseCase{
		apiKeyRepo: apiKeyRepo,
	}
}
