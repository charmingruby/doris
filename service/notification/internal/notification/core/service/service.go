package service

import (
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/repository"
)

type Service struct {
	logger           *instrumentation.Logger
	notificationRepo repository.NotificationRepository
}

func New(logger *instrumentation.Logger, notificationRepo repository.NotificationRepository) *Service {
	return &Service{
		logger:           logger,
		notificationRepo: notificationRepo,
	}
}
