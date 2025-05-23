package usecase

import (
	"context"
	"errors"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/model"
)

func (s *Suite) Test_DispatchNotification() {
	validInput := DispatchNotificationInput{
		CorrelationID:    "123",
		To:               "test@test.com",
		RecipientName:    "Test",
		Content:          "Test",
		NotificationType: model.OTPNotification,
	}

	s.Run("it should dispatch a notification", func() {
		err := s.uc.DispatchNotification(context.Background(), validInput)
		s.NoError(err)

		storedNotification := s.notificationRepo.Items[0]
		sentNotification := s.notifierClient.Items[0]

		s.Equal(1, len(s.notificationRepo.Items))
		s.Equal(validInput.CorrelationID, storedNotification.CorrelationID)
		s.Equal(validInput.To, storedNotification.To)
		s.Equal(validInput.RecipientName, storedNotification.RecipientName)
		s.Equal(validInput.NotificationType, storedNotification.NotificationType)

		s.Equal(1, len(s.notifierClient.Items))
		s.Equal(validInput.CorrelationID, sentNotification.CorrelationID)
		s.Equal(validInput.To, sentNotification.To)
		s.Equal(validInput.RecipientName, sentNotification.RecipientName)
		s.Equal(validInput.NotificationType, sentNotification.NotificationType)
	})

	s.Run("it should be not able to dispatch a notification if the notification client fails", func() {
		s.notifierClient.IsHealthy = false

		err := s.uc.DispatchNotification(context.Background(), validInput)

		s.Error(err)
		var errExternalService *custom_err.ErrExternalService
		s.True(errors.As(err, &errExternalService), "error should be of type ErrExternalService")

		s.Equal(0, len(s.notificationRepo.Items))
		s.Equal(0, len(s.notifierClient.Items))
	})

	s.Run("it should be not able to dispatch a notification if the notification storage fails", func() {
		s.notificationRepo.IsHealthy = false

		err := s.uc.DispatchNotification(context.Background(), validInput)
		s.Error(err)

		var errDatasourceOperationFailed *custom_err.ErrDatasourceOperationFailed
		s.True(errors.As(err, &errDatasourceOperationFailed), "error should be of type ErrDatasourceOperationFailed")

		s.Equal(0, len(s.notificationRepo.Items))
		s.Equal(1, len(s.notifierClient.Items))
	})
}
