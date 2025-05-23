package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/id"
	"github.com/charmingruby/doris/service/account/internal/access/core/event"
	"github.com/charmingruby/doris/service/account/internal/access/core/model"
)

func (s *Suite) Test_ActivateAPIKey() {
	expirationDelay := 10 * time.Minute

	dummyAPIKey := *model.NewAPIKey(model.APIKeyInput{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Key:       id.New(),
	})

	dummyOTP, err := model.NewOTP(model.OTPInput{
		Purpose:       model.OTP_PURPOSE_API_KEY_ACTIVATION,
		CorrelationID: dummyAPIKey.ID,
		ExpiresAt:     time.Now().Add(expirationDelay),
	})

	s.NoError(err)

	s.Run("it should activate the api key", func() {
		ctx := context.Background()

		err := s.apiKeyRepo.Create(ctx, dummyAPIKey)
		s.NoError(err)

		err = s.otpRepo.Create(ctx, *dummyOTP)
		s.NoError(err)

		s.Equal(0, len(s.evtHandler.Pub.Messages))

		err = s.evtHandler.DispatchSendOTPNotification(ctx, event.SendOTPNotification{
			ID:            dummyAPIKey.ID,
			To:            dummyAPIKey.Email,
			RecipientName: dummyAPIKey.FirstName + " " + dummyAPIKey.LastName,
			Code:          dummyOTP.Code,
			SentAt:        time.Now(),
		})
		s.NoError(err)

		s.Equal(1, len(s.evtHandler.Pub.Messages))

		token, err := s.uc.ActivateAPIKey(ctx, ActivateAPIKeyInput{
			APIKeyID: dummyAPIKey.ID,
			OTP:      dummyOTP.Code,
		})
		s.NoError(err)
		s.NotEmpty(token)

		verifiedAPIKey := s.apiKeyRepo.Items[0]
		s.Equal(model.API_KEY_STATUS_ACTIVE, verifiedAPIKey.Status)

		generatedToken := verifiedAPIKey.ID + "-token"
		tokenPayload := s.tokenClient.Items[generatedToken]
		s.Equal(verifiedAPIKey.ID, tokenPayload.Sub)
		s.Equal(verifiedAPIKey.Tier, tokenPayload.Payload.Tier)
	})

	s.Run("it should be not able to activate the api key if the datasource operation fails", func() {
		ctx := context.Background()

		s.apiKeyRepo.IsHealthy = false

		_, err := s.uc.ActivateAPIKey(ctx, ActivateAPIKeyInput{
			APIKeyID: dummyAPIKey.ID,
			OTP:      dummyOTP.Code,
		})

		s.Error(err)

		var dsErr *custom_err.ErrDatasourceOperationFailed
		s.True(errors.As(err, &dsErr), "error should be of type ErrDatasourceOperationFailed")
	})

	s.Run("it should be not able to activate the api key if the api key is not found", func() {
		ctx := context.Background()

		_, err := s.uc.ActivateAPIKey(ctx, ActivateAPIKeyInput{
			APIKeyID: dummyAPIKey.ID,
			OTP:      dummyOTP.Code,
		})

		s.Error(err)

		var resourceNotFoundErr *custom_err.ErrResourceNotFound
		s.True(errors.As(err, &resourceNotFoundErr), "error should be of type ErrResourceNotFound")
	})

	s.Run("it should be not able to activate the api key if the otp is not found", func() {
		ctx := context.Background()

		err := s.apiKeyRepo.Create(ctx, dummyAPIKey)
		s.NoError(err)

		_, err = s.uc.ActivateAPIKey(ctx, ActivateAPIKeyInput{
			APIKeyID: dummyAPIKey.ID,
			OTP:      dummyOTP.Code,
		})

		s.Error(err)

		var resourceNotFoundErr *custom_err.ErrResourceNotFound
		s.True(errors.As(err, &resourceNotFoundErr), "error should be of type ErrResourceNotFound")
	})

	s.Run("it should be not able to activate the api key if the confirmation code does not match", func() {
		ctx := context.Background()

		err := s.apiKeyRepo.Create(ctx, dummyAPIKey)
		s.NoError(err)

		err = s.otpRepo.Create(ctx, *dummyOTP)
		s.NoError(err)

		s.Equal(0, len(s.evtHandler.Pub.Messages))

		err = s.evtHandler.DispatchSendOTPNotification(ctx, event.SendOTPNotification{
			ID:            dummyAPIKey.ID,
			To:            dummyAPIKey.Email,
			RecipientName: dummyAPIKey.FirstName + " " + dummyAPIKey.LastName,
			Code:          dummyOTP.Code,
			SentAt:        time.Now(),
		})
		s.NoError(err)

		s.Equal(1, len(s.evtHandler.Pub.Messages))

		_, err = s.uc.ActivateAPIKey(ctx, ActivateAPIKeyInput{
			APIKeyID: dummyAPIKey.ID,
			OTP:      "invalid-code",
		})

		s.Error(err)

		var invalidOTPCodeErr *custom_err.ErrInvalidOTPCode
		s.True(errors.As(err, &invalidOTPCodeErr), "error should be of type ErrInvalidOTPCode")
	})

	s.Run("it should be not able to activate the api key if the confirmation code has expired", func() {
		ctx := context.Background()

		expiredOTP := *dummyOTP
		expiredOTP.ExpiresAt = time.Now().Add(-1 * time.Hour)

		err := s.otpRepo.Create(ctx, expiredOTP)
		s.NoError(err)

		err = s.apiKeyRepo.Create(ctx, dummyAPIKey)
		s.NoError(err)

		s.Equal(0, len(s.evtHandler.Pub.Messages))

		err = s.evtHandler.DispatchSendOTPNotification(ctx, event.SendOTPNotification{
			ID:            dummyAPIKey.ID,
			To:            dummyAPIKey.Email,
			RecipientName: dummyAPIKey.FirstName + " " + dummyAPIKey.LastName,
			Code:          expiredOTP.Code,
			SentAt:        time.Now(),
		})
		s.NoError(err)

		s.Equal(1, len(s.evtHandler.Pub.Messages))

		_, err = s.uc.ActivateAPIKey(ctx, ActivateAPIKeyInput{
			APIKeyID: dummyAPIKey.ID,
			OTP:      expiredOTP.Code,
		})

		s.Error(err)

		var invalidOTPCodeErr *custom_err.ErrInvalidOTPCode
		s.True(errors.As(err, &invalidOTPCodeErr), "error should be of type ErrInvalidOTPCode")
	})

	s.Run("it should be not able to activate the api key if the api key is already activated", func() {
		ctx := context.Background()

		err := s.apiKeyRepo.Create(ctx, dummyAPIKey)
		s.NoError(err)

		err = s.otpRepo.Create(ctx, *dummyOTP)
		s.NoError(err)

		s.Equal(0, len(s.evtHandler.Pub.Messages))

		err = s.evtHandler.DispatchSendOTPNotification(ctx, event.SendOTPNotification{
			ID:            dummyAPIKey.ID,
			To:            dummyAPIKey.Email,
			RecipientName: dummyAPIKey.FirstName + " " + dummyAPIKey.LastName,
			Code:          dummyOTP.Code,
			SentAt:        time.Now(),
		})
		s.NoError(err)

		s.Equal(1, len(s.evtHandler.Pub.Messages))

		_, err = s.uc.ActivateAPIKey(ctx, ActivateAPIKeyInput{
			APIKeyID: dummyAPIKey.ID,
			OTP:      dummyOTP.Code,
		})
		s.NoError(err)

		_, err = s.uc.ActivateAPIKey(ctx, ActivateAPIKeyInput{
			APIKeyID: dummyAPIKey.ID,
			OTP:      dummyOTP.Code,
		})
		s.Error(err)

		var apiKeyAlreadyActivatedErr *custom_err.ErrAPIKeyAlreadyActivated
		s.True(errors.As(err, &apiKeyAlreadyActivatedErr), "error should be of type ErrAPIKeyAlreadyActivated")
	})
}
