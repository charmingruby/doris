package client

import (
	"context"

	"github.com/charmingruby/doris/service/notification/internal/notification/core/model"
)

type Notifier interface {
	Send(ctx context.Context, notification model.Notification) error
}
