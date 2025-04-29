package memory

import (
	"context"

	"github.com/charmingruby/doris/service/notification/internal/notification/core/model"
)

type Notifier struct {
	Items     []model.Notification
	IsHealthy bool
}

func NewNotifier() *Notifier {
	return &Notifier{
		Items:     []model.Notification{},
		IsHealthy: true,
	}
}

func (n *Notifier) Send(ctx context.Context, notification model.Notification) error {
	if !n.IsHealthy {
		return ErrUnhealthyDatasource
	}

	n.Items = append(n.Items, notification)

	return nil
}
