package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/id"
	"github.com/charmingruby/doris/service/account/internal/access/core/model"
)

func (s *Suite) Test_ResendAPIKeyActivation() {
	validInput := ResendAPIKeyActivationInput{
		APIKeyID: id.New(),
	}

	dummyAPIKey := *model.NewAPIKey(model.APIKeyInput{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Key:       id.New(),
	})
	dummyAPIKey.ID = validInput.APIKeyID

	s.Run("it should resend api key activation", func() {
		err := s.apiKeyRepo.Create(context.Background(), dummyAPIKey)
		s.NoError(err)

		err = s.uc.ResendAPIKeyActivation(context.Background(), validInput)
		s.NoError(err)

		otp := s.otpRepo.Items[0]
		expectedExpiration := time.Now().Add(30 * time.Minute)

		timeDiff := otp.ExpiresAt.Sub(expectedExpiration)
		s.True(timeDiff < time.Second && timeDiff > -time.Second, "expiration time should be within 1 second of expected time")

		s.Equal(1, len(s.evtHandler.Pub.Messages))
	})

	s.Run("it should not be able to resend if api key is not found", func() {
		err := s.uc.ResendAPIKeyActivation(context.Background(), validInput)
		s.Error(err)

		var errResourceNotFound *custom_err.ErrResourceNotFound
		s.True(errors.As(err, &errResourceNotFound), "error should be of type ErrResourceNotFound")
	})

	s.Run("it should not be able to resend if datasource fails", func() {
		err := s.apiKeyRepo.Create(context.Background(), dummyAPIKey)
		s.NoError(err)

		s.apiKeyRepo.IsHealthy = false

		err = s.uc.ResendAPIKeyActivation(context.Background(), validInput)
		s.Error(err)

		var dsErr *custom_err.ErrDatasourceOperationFailed
		s.True(errors.As(err, &dsErr), "error should be of type ErrDatasourceOperationFailed")
	})

	s.Run("it should not be able to resend if messaging fails", func() {
		err := s.apiKeyRepo.Create(context.Background(), dummyAPIKey)
		s.NoError(err)

		s.evtHandler.Pub.IsHealthy = false

		err = s.uc.ResendAPIKeyActivation(context.Background(), validInput)
		s.Error(err)

		var errMessaging *custom_err.ErrMessagingWrapper
		s.True(errors.As(err, &errMessaging), "error should be of type ErrMessagingWrapper")
	})

	s.Run("it should not be able to resend if there is a cooldown period", func() {
		err := s.apiKeyRepo.Create(context.Background(), dummyAPIKey)
		s.NoError(err)

		otp, err := model.NewOTP(model.OTPInput{
			Purpose:       model.OTP_PURPOSE_API_KEY_ACTIVATION,
			CorrelationID: dummyAPIKey.ID,
			ExpiresAt:     time.Now().UTC().Add(30 * time.Minute),
		})
		s.NoError(err)

		err = s.otpRepo.Create(context.Background(), *otp)
		s.NoError(err)

		err = s.uc.ResendAPIKeyActivation(context.Background(), validInput)
		s.Error(err)

		var errOTPCooldown *custom_err.ErrOTPGenerationCooldown
		s.True(errors.As(err, &errOTPCooldown), "error should be of type ErrOTPGenerationCooldown")
	})
}
