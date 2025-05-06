package memory

import (
	"context"

	"github.com/charmingruby/doris/lib/core/pagination"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/model"
)

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

func (r *NotificationRepository) FindManyByCorrelationID(ctx context.Context, correlationID string, page int) ([]model.Notification, error) {
	if !r.IsHealthy {
		return nil, ErrUnhealthyDatasource
	}

	filteredItems := []model.Notification{}
	for _, item := range r.Items {
		if item.CorrelationID == correlationID {
			filteredItems = append(filteredItems, item)
		}
	}

	startIdx := (page - 1) * pagination.MAX_ITEMS_PER_PAGE
	endIdx := startIdx + pagination.MAX_ITEMS_PER_PAGE

	if startIdx >= len(filteredItems) {
		return []model.Notification{}, nil
	}

	if endIdx > len(filteredItems) {
		endIdx = len(filteredItems)
	}

	return filteredItems[startIdx:endIdx], nil
}
