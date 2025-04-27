package service

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/model"
)

type NotifyApiKeyActivationInput struct {
	CorrelationID string
	To            string
	RecipientName string
	EmittedAt     time.Time
}

func (s *Service) NotifyApiKeyActivation(ctx context.Context, in NotifyApiKeyActivationInput) error {
	notification := model.NewNotification(model.NotificationInput{
		CorrelationID: in.CorrelationID,
		To:            in.To,
		RecipientName: in.RecipientName,
		MessageType:   model.APIKeyActivation,
		EmittedAt:     in.EmittedAt,
	})

	if err := s.notificationRepo.Create(ctx, *notification); err != nil {
		return custom_err.NewErrDatasourceOperationFailed("create notification", err)
	}

	// notify the recipient

	return nil
}
