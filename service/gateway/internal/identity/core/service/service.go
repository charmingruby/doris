package service

import (
	"github.com/charmingruby/doris/service/gateway/internal/identity/core/repository"
)

type Service struct {
	apiKeyRepo repository.APIKeyRepository
}

func New(
	apiKeyRepo repository.APIKeyRepository,

) *Service {
	return &Service{
		apiKeyRepo: apiKeyRepo,
	}
}
