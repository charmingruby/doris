package repository

import (
	"context"

	"github.com/charmingruby/doris/service/notification/internal/notification/core/model"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification model.Notification) error
}
