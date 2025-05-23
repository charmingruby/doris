package usecase

import (
	"context"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/model"
)

type DispatchNotificationInput struct {
	CorrelationID    string
	To               string
	RecipientName    string
	Content          string
	NotificationType model.NotificationType
}

func (uc *UseCase) DispatchNotification(ctx context.Context, in DispatchNotificationInput) error {
	notification := model.NewNotification(model.NotificationInput{
		CorrelationID:    in.CorrelationID,
		To:               in.To,
		RecipientName:    in.RecipientName,
		Content:          in.Content,
		NotificationType: in.NotificationType,
	})

	if err := uc.notifier.Send(ctx, *notification); err != nil {
		return custom_err.NewErrExternalService(err)
	}

	if err := uc.repo.Create(ctx, *notification); err != nil {
		return custom_err.NewErrDatasourceOperationFailed("create notification", err)
	}

	return nil
}
