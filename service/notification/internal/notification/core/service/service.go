package service

import (
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/repository"
)

type Service struct {
	logger *instrumentation.Logger
	repo   repository.NotificationRepository
}

func New(logger *instrumentation.Logger, repo repository.NotificationRepository) *Service {
	return &Service{
		logger: logger,
		repo:   repo,
	}
}
