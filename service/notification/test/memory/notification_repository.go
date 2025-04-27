package memory

import (
	"context"
	"errors"

	"github.com/charmingruby/doris/service/notification/internal/notification/core/model"
)

var ErrUnhealthyDatasource = errors.New("datasource is unhealthy")

type NotificationRepository struct {
	Items     []model.Notification
	IsHealthy bool
}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{
		Items:     []model.Notification{},
		IsHealthy: true,
	}
}

func (r *NotificationRepository) Create(ctx context.Context, notification model.Notification) error {
	if !r.IsHealthy {
		return ErrUnhealthyDatasource
	}

	r.Items = append(r.Items, notification)

	return nil
}
