package service

import (
	"context"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/model"
)

type NotifyOTPInput struct {
	CorrelationID string
	To            string
	RecipientName string
	Content       string
}

func (s *Service) NotifyOTP(ctx context.Context, in NotifyOTPInput) error {
	notification := model.NewNotification(model.NotificationInput{
		CorrelationID: in.CorrelationID,
		To:            in.To,
		RecipientName: in.RecipientName,
		Content:       in.Content,
		MessageType:   model.OTP,
	})

	if err := s.repo.Create(ctx, *notification); err != nil {
		return custom_err.NewErrDatasourceOperationFailed("create notification", err)
	}

	// notify the recipient

	return nil
}
