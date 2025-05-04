package service

import (
	"context"
	"errors"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/id"
	"github.com/charmingruby/doris/service/account/internal/access/core/event"
	"github.com/charmingruby/doris/service/account/internal/access/core/model"
)

func (s *Suite) Test_VerifySignInIntent() {
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

	s.Run("it should verify a sign in intent with valid otp", func() {
		ctx := context.Background()

		validAPIKey := dummyAPIKey
		validAPIKey.Status = model.API_KEY_STATUS_ACTIVE

		err := s.apiKeyRepo.Create(ctx, validAPIKey)
		s.NoError(err)

		err = s.otpRepo.Create(ctx, *dummyOTP)
		s.NoError(err)

		s.Equal(0, len(s.evtHandler.Pub.Messages))

		err = s.evtHandler.DispatchSendOTPNotification(ctx, event.SendOTPNotificationMessage{
			ID:            validAPIKey.ID,
			To:            validAPIKey.Email,
			RecipientName: validAPIKey.FirstName + " " + validAPIKey.LastName,
			Code:          dummyOTP.Code,
			SentAt:        time.Now(),
		})
		s.NoError(err)

		s.Equal(1, len(s.evtHandler.Pub.Messages))

		token, err := s.svc.VerifySignInIntent(ctx, VerifySignInIntentInput{
			Email: validAPIKey.Email,
			OTP:   dummyOTP.Code,
		})

		s.NoError(err)
		s.NotEmpty(token)

		generatedToken := validAPIKey.ID + "-token"
		tokenPayload := s.tokenClient.Items[generatedToken]
		s.Equal(validAPIKey.ID, tokenPayload.Sub)
		s.Equal(validAPIKey.Tier, tokenPayload.Payload.Tier)
	})

	s.Run("it should be not able to verify sign in intent if the datasource operation fails", func() {
		ctx := context.Background()

		validAPIKey := dummyAPIKey
		validAPIKey.Status = model.API_KEY_STATUS_ACTIVE

		err := s.apiKeyRepo.Create(ctx, validAPIKey)
		s.NoError(err)

		err = s.otpRepo.Create(ctx, *dummyOTP)
		s.NoError(err)

		s.Equal(0, len(s.evtHandler.Pub.Messages))

		err = s.evtHandler.DispatchSendOTPNotification(ctx, event.SendOTPNotificationMessage{
			ID:            validAPIKey.ID,
			To:            validAPIKey.Email,
			RecipientName: validAPIKey.FirstName + " " + validAPIKey.LastName,
			Code:          dummyOTP.Code,
			SentAt:        time.Now(),
		})
		s.NoError(err)

		s.Equal(1, len(s.evtHandler.Pub.Messages))

		s.apiKeyRepo.IsHealthy = false

		token, err := s.svc.VerifySignInIntent(ctx, VerifySignInIntentInput{
			Email: validAPIKey.Email,
			OTP:   dummyOTP.Code,
		})

		s.Error(err)
		s.Empty(token)

		var dsErr *custom_err.ErrDatasourceOperationFailed
		s.True(errors.As(err, &dsErr), "error should be of type ErrDatasourceOperationFailed")
	})

	s.Run("it should be not able verify a sign in intent if the api key is not found", func() {
		ctx := context.Background()

		_, err := s.svc.VerifySignInIntent(ctx, VerifySignInIntentInput{
			Email: "invalid email",
			OTP:   dummyOTP.Code,
		})

		s.Error(err)

		var resourceNotFoundErr *custom_err.ErrResourceNotFound
		s.True(errors.As(err, &resourceNotFoundErr), "error should be of type ErrResourceNotFound")
	})

	s.Run("it should be not able to verify a sign in intent if the otp is not found", func() {
		ctx := context.Background()

		validAPIKey := dummyAPIKey
		validAPIKey.Status = model.API_KEY_STATUS_ACTIVE

		err := s.apiKeyRepo.Create(ctx, validAPIKey)
		s.NoError(err)

		_, err = s.svc.VerifySignInIntent(ctx, VerifySignInIntentInput{
			Email: validAPIKey.Email,
			OTP:   dummyOTP.Code,
		})

		s.Error(err)

		var resourceNotFoundErr *custom_err.ErrResourceNotFound
		s.True(errors.As(err, &resourceNotFoundErr), "error should be of type ErrResourceNotFound")
	})

	s.Run("it should be not able to verify a sign in intent if the confirmation code does not match", func() {
		ctx := context.Background()

		validAPIKey := dummyAPIKey
		validAPIKey.Status = model.API_KEY_STATUS_ACTIVE

		err := s.apiKeyRepo.Create(ctx, validAPIKey)
		s.NoError(err)

		err = s.otpRepo.Create(ctx, *dummyOTP)
		s.NoError(err)

		s.Equal(0, len(s.evtHandler.Pub.Messages))

		err = s.evtHandler.DispatchSendOTPNotification(ctx, event.SendOTPNotificationMessage{
			ID:            validAPIKey.ID,
			To:            validAPIKey.Email,
			RecipientName: validAPIKey.FirstName + " " + validAPIKey.LastName,
			Code:          dummyOTP.Code,
			SentAt:        time.Now(),
		})
		s.NoError(err)

		s.Equal(1, len(s.evtHandler.Pub.Messages))

		_, err = s.svc.VerifySignInIntent(ctx, VerifySignInIntentInput{
			Email: validAPIKey.Email,
			OTP:   "invalid-code",
		})

		s.Error(err)

		var invalidOTPCodeErr *custom_err.ErrInvalidOTPCode
		s.True(errors.As(err, &invalidOTPCodeErr), "error should be of type ErrInvalidOTPCode")
	})

	s.Run("it should be not able to verify a sign in intent if the confirmation code has expired", func() {
		ctx := context.Background()

		expiredOTP := *dummyOTP
		expiredOTP.ExpiresAt = time.Now().Add(-1 * time.Hour)

		validAPIKey := dummyAPIKey
		validAPIKey.Status = model.API_KEY_STATUS_ACTIVE

		err = s.apiKeyRepo.Create(ctx, validAPIKey)
		s.NoError(err)

		err = s.otpRepo.Create(ctx, expiredOTP)
		s.NoError(err)

		s.Equal(0, len(s.evtHandler.Pub.Messages))

		err = s.evtHandler.DispatchSendOTPNotification(ctx, event.SendOTPNotificationMessage{
			ID:            validAPIKey.ID,
			To:            validAPIKey.Email,
			RecipientName: validAPIKey.FirstName + " " + validAPIKey.LastName,
			Code:          expiredOTP.Code,
			SentAt:        time.Now(),
		})
		s.NoError(err)

		s.Equal(1, len(s.evtHandler.Pub.Messages))

		_, err = s.svc.VerifySignInIntent(ctx, VerifySignInIntentInput{
			Email: validAPIKey.Email,
			OTP:   expiredOTP.Code,
		})

		s.Error(err)

		var invalidOTPCodeErr *custom_err.ErrInvalidOTPCode
		s.True(errors.As(err, &invalidOTPCodeErr), "error should be of type ErrInvalidOTPCode")
	})

	s.Run("it should be not able to verify a sign in intent if api key has insufficient permission", func() {
		invalidAPIKey := dummyAPIKey
		invalidAPIKey.Status = model.API_KEY_STATUS_PENDING

		err := s.apiKeyRepo.Create(context.Background(), invalidAPIKey)
		s.NoError(err)

		ctx := context.Background()

		validAPIKey := dummyAPIKey
		validAPIKey.Status = model.API_KEY_STATUS_ACTIVE

		err = s.apiKeyRepo.Create(ctx, validAPIKey)
		s.NoError(err)

		err = s.otpRepo.Create(ctx, *dummyOTP)
		s.NoError(err)

		s.Equal(0, len(s.evtHandler.Pub.Messages))

		err = s.evtHandler.DispatchSendOTPNotification(ctx, event.SendOTPNotificationMessage{
			ID:            validAPIKey.ID,
			To:            validAPIKey.Email,
			RecipientName: validAPIKey.FirstName + " " + validAPIKey.LastName,
			Code:          dummyOTP.Code,
			SentAt:        time.Now(),
		})
		s.NoError(err)

		s.Equal(1, len(s.evtHandler.Pub.Messages))

		token, err := s.svc.VerifySignInIntent(ctx, VerifySignInIntentInput{
			Email: validAPIKey.Email,
			OTP:   dummyOTP.Code,
		})

		s.Error(err)
		s.Empty(token)

		var errInsufficientPermission *custom_err.ErrInsufficientPermission
		s.True(errors.As(err, &errInsufficientPermission), "error should be of type ErrInsufficientPermission")
	})
}
