package service

import (
	"context"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/model"
)

type FetchCorrelatedNotificationsInput struct {
	CorrelationID string `json:"correlation_id"`
	Page          int    `json:"page"`
}

type FetchCorrelatedNotificationsOutput struct {
	Notifications []model.Notification `json:"notifications"`
}

func (s *Service) FetchCorrelatedNotifications(ctx context.Context, in FetchCorrelatedNotificationsInput) (FetchCorrelatedNotificationsOutput, error) {
	notifications, err := s.repo.FindManyByCorrelationID(ctx, in.CorrelationID, in.Page)
	if err != nil {
		return FetchCorrelatedNotificationsOutput{}, custom_err.NewErrDatasourceOperationFailed("find many notifications by correlation_id", err)
	}

	return FetchCorrelatedNotificationsOutput{
		Notifications: notifications,
	}, nil
}
