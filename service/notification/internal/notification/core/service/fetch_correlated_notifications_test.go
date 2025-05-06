package service

import (
	"context"
	"errors"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/pagination"
	"github.com/charmingruby/doris/service/notification/internal/notification/core/model"
)

func (s *Suite) Test_FetchCorrelatedNotifications() {
	validInput := FetchCorrelatedNotificationsInput{
		CorrelationID: "test-correlation-id",
		Page:          1,
	}

	dummyNotification := model.Notification{
		ID:               "test-id",
		CorrelationID:    validInput.CorrelationID,
		To:               "test@example.com",
		RecipientName:    "Test User",
		Content:          "Test Content",
		NotificationType: model.OTPNotification,
		CreatedAt:        time.Now(),
	}

	s.Run("it should fetch notifications by correlation id", func() {
		err := s.notificationRepo.Create(context.Background(), dummyNotification)
		s.NoError(err)

		output, err := s.svc.FetchCorrelatedNotifications(context.Background(), validInput)
		s.NoError(err)
		s.Equal(1, len(output.Notifications))
		s.Equal(dummyNotification.ID, output.Notifications[0].ID)
		s.Equal(dummyNotification.CorrelationID, output.Notifications[0].CorrelationID)
		s.Equal(dummyNotification.To, output.Notifications[0].To)
		s.Equal(dummyNotification.RecipientName, output.Notifications[0].RecipientName)
		s.Equal(dummyNotification.Content, output.Notifications[0].Content)
		s.Equal(dummyNotification.NotificationType, output.Notifications[0].NotificationType)
	})

	s.Run("it should return empty list when no notifications found", func() {
		output, err := s.svc.FetchCorrelatedNotifications(context.Background(), validInput)
		s.NoError(err)
		s.Equal(0, len(output.Notifications))
	})

	s.Run("it should not be able to fetch notifications if datasource fails", func() {
		s.notificationRepo.IsHealthy = false

		output, err := s.svc.FetchCorrelatedNotifications(context.Background(), validInput)
		s.Error(err)
		s.Equal(0, len(output.Notifications))

		var dsErr *custom_err.ErrDatasourceOperationFailed
		s.True(errors.As(err, &dsErr), "error should be of type ErrDatasourceOperationFailed")
	})

	s.Run("it should handle pagination correctly", func() {
		for i := range pagination.MAX_ITEMS_PER_PAGE + 5 {
			notification := model.Notification{
				ID:               "test-id-" + string(rune(i)),
				CorrelationID:    validInput.CorrelationID,
				To:               "test@example.com",
				RecipientName:    "Test User",
				Content:          "Test Content",
				NotificationType: model.OTPNotification,
				CreatedAt:        time.Now(),
			}
			err := s.notificationRepo.Create(context.Background(), notification)
			s.NoError(err)
		}

		output, err := s.svc.FetchCorrelatedNotifications(context.Background(), validInput)
		s.NoError(err)
		s.Equal(pagination.MAX_ITEMS_PER_PAGE, len(output.Notifications))

		validInput.Page = 2
		output, err = s.svc.FetchCorrelatedNotifications(context.Background(), validInput)
		s.NoError(err)
		s.Equal(5, len(output.Notifications))
	})
}
