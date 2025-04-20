package usecase

import (
	"github.com/charmingruby/doris/service/gateway/internal/identity/core/client"
	"github.com/charmingruby/doris/service/gateway/internal/identity/core/repository"
)

type UseCase struct {
	apiKeyRepo  repository.APIKeyRepository
	emailClient client.EmailClient
}

func New(
	apiKeyRepo repository.APIKeyRepository,
	emailClient client.EmailClient,
) *UseCase {
	return &UseCase{
		apiKeyRepo:  apiKeyRepo,
		emailClient: emailClient,
	}
}
