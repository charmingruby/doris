package repository

import (
	"context"

	"github.com/charmingruby/doris/service/notification/internal/notification/core/model"
)

type NotificationRepository interface {
	FindManyByCorrelationID(ctx context.Context, correlationID string, page int) ([]model.Notification, error)
	Create(ctx context.Context, notification model.Notification) error
}
