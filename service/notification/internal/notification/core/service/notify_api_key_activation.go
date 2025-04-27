package service

import (
	"context"
	"time"

	"github.com/charmingruby/doris/service/notification/internal/notification/core/model"
)

type NotifyApiKeyActivationInput struct {
	CorrelationID string
	To            string
	RecipientName string
}

func (s *Service) NotifyApiKeyActivation(ctx context.Context, in NotifyApiKeyActivationInput) error {
	notification := model.NewNotification(model.NotificationInput{
		CorrelationID: in.CorrelationID,
		To:            in.To,
		RecipientName: in.RecipientName,
		MessageType:   model.APIKeyActivation,
		EmittedAt:     time.Now(),
	})

	s.logger.Info("notification created", "notification", notification)

	return nil
}
