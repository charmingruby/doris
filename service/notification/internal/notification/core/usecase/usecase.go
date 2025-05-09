package usecase

import (
	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/client"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/repository"
)

type UseCase struct {
	logger   *instrumentation.Logger
	repo     repository.NotificationRepository
	notifier client.Notifier
}

func New(logger *instrumentation.Logger, repo repository.NotificationRepository, notifier client.Notifier) *UseCase {
	return &UseCase{
		logger:   logger,
		repo:     repo,
		notifier: notifier,
	}
}
