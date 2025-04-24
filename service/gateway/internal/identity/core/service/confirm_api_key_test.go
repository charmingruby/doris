package service

import (
	"context"
	"errors"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/id"
	"github.com/charmingruby/doris/service/gateway/internal/identity/core/event"
	"github.com/charmingruby/doris/service/gateway/internal/identity/core/model"
)

func (s *Suite) Test_ConfirmAPIKey() {
	expirationDelay := 10 * time.Minute

	dummyAPIKey := *model.NewAPIKey(model.APIKeyInput{
		FirstName:                 "John",
		LastName:                  "Doe",
		Email:                     "john.doe@example.com",
		Key:                       id.New(),
		ConfirmationCodeExpiresAt: time.Now().Add(expirationDelay),
	})

	s.Run("it should confirm the api key", func() {
		ctx := context.Background()

		err := s.apiKeyRepo.Create(ctx, dummyAPIKey)
		s.NoError(err)

		s.Equal(0, len(s.pub.Messages))

		err = s.evtHandler.PublishAPIKeyRequest(ctx, &event.APIKeyRequest{
			ID:               dummyAPIKey.ID,
			To:               dummyAPIKey.Email,
			ConfirmationCode: dummyAPIKey.ConfirmationCode,
		})
		s.NoError(err)

		s.Equal(1, len(s.pub.Messages))

		storedApiKey := s.apiKeyRepo.Items[0]

		err = s.svc.ConfirmAPIKey(ctx, ConfirmAPIKeyInput{
			Key:              dummyAPIKey.Key,
			ConfirmationCode: storedApiKey.ConfirmationCode,
		})

		s.NoError(err)

		verifiedAPIKey := s.apiKeyRepo.Items[0]

		s.Equal(model.API_KEY_STATUS_ACTIVE, verifiedAPIKey.Status)
	})

	s.Run("it should be not able to confirm the api key if the datasource operation fails", func() {
		ctx := context.Background()

		s.apiKeyRepo.IsHealthy = false

		err := s.svc.ConfirmAPIKey(ctx, ConfirmAPIKeyInput{
			Key:              dummyAPIKey.Key,
			ConfirmationCode: dummyAPIKey.ConfirmationCode,
		})

		s.Error(err)

		var dsErr *custom_err.ErrDatasourceOperationFailed
		s.True(errors.As(err, &dsErr), "error should be of type ErrDatasourceOperationFailed")
	})

	s.Run("it should be not able to confirm the api key if the api key is not found", func() {
		ctx := context.Background()

		err := s.svc.ConfirmAPIKey(ctx, ConfirmAPIKeyInput{
			Key:              dummyAPIKey.Key,
			ConfirmationCode: dummyAPIKey.ConfirmationCode,
		})

		s.Error(err)

		var resourceNotFoundErr *custom_err.ErrResourceNotFound
		s.True(errors.As(err, &resourceNotFoundErr), "error should be of type ErrResourceNotFound")
	})

	s.Run("it should be not able to confirm the api key if the confirmation code does not match", func() {
		ctx := context.Background()

		err := s.apiKeyRepo.Create(ctx, dummyAPIKey)
		s.NoError(err)

		s.Equal(0, len(s.pub.Messages))

		err = s.evtHandler.PublishAPIKeyRequest(ctx, &event.APIKeyRequest{
			ID:               dummyAPIKey.ID,
			To:               dummyAPIKey.Email,
			ConfirmationCode: dummyAPIKey.ConfirmationCode,
		})
		s.NoError(err)

		s.Equal(1, len(s.pub.Messages))

		err = s.svc.ConfirmAPIKey(ctx, ConfirmAPIKeyInput{
			Key:              dummyAPIKey.Key,
			ConfirmationCode: "invalid-code",
		})

		s.Error(err)

		var invalidConfirmationCodeErr *custom_err.ErrInvalidConfirmationCode
		s.True(errors.As(err, &invalidConfirmationCodeErr), "error should be of type ErrInvalidConfirmationCode")
	})

	s.Run("it should be not able to confirm the api key if the confirmation code has expired", func() {
		ctx := context.Background()

		dummyAPIKeyClone := dummyAPIKey

		alwaysExpiredDate := time.Date(1920, 1, 1, 0, 0, 0, 0, time.UTC)
		dummyAPIKeyClone.ConfirmationCodeExpiresAt = alwaysExpiredDate

		err := s.apiKeyRepo.Create(ctx, dummyAPIKeyClone)
		s.NoError(err)

		s.Equal(0, len(s.pub.Messages))

		err = s.evtHandler.PublishAPIKeyRequest(ctx, &event.APIKeyRequest{
			ID:               dummyAPIKeyClone.ID,
			To:               dummyAPIKeyClone.Email,
			ConfirmationCode: dummyAPIKeyClone.ConfirmationCode,
		})
		s.NoError(err)

		s.Equal(1, len(s.pub.Messages))

		err = s.svc.ConfirmAPIKey(ctx, ConfirmAPIKeyInput{
			Key:              dummyAPIKeyClone.Key,
			ConfirmationCode: "invalid-code",
		})

		s.Error(err)

		var invalidConfirmationCodeErr *custom_err.ErrInvalidConfirmationCode
		s.True(errors.As(err, &invalidConfirmationCodeErr), "error should be of type ErrInvalidConfirmationCode")
	})

	s.Run("it should be not able to confirm the api key if the api key is already confirmed", func() {
		ctx := context.Background()

		err := s.apiKeyRepo.Create(ctx, dummyAPIKey)
		s.NoError(err)

		s.Equal(0, len(s.pub.Messages))

		err = s.evtHandler.PublishAPIKeyRequest(ctx, &event.APIKeyRequest{
			ID:               dummyAPIKey.ID,
			To:               dummyAPIKey.Email,
			ConfirmationCode: dummyAPIKey.ConfirmationCode,
		})
		s.NoError(err)

		s.Equal(1, len(s.pub.Messages))

		storedApiKey := s.apiKeyRepo.Items[0]

		err = s.svc.ConfirmAPIKey(ctx, ConfirmAPIKeyInput{
			Key:              dummyAPIKey.Key,
			ConfirmationCode: storedApiKey.ConfirmationCode,
		})

		s.NoError(err)

		err = s.svc.ConfirmAPIKey(ctx, ConfirmAPIKeyInput{
			Key:              dummyAPIKey.Key,
			ConfirmationCode: storedApiKey.ConfirmationCode,
		})

		s.Error(err)

		var apiKeyAlreadyConfirmedErr *custom_err.ErrAPIKeyAlreadyConfirmed
		s.True(errors.As(err, &apiKeyAlreadyConfirmedErr), "error should be of type ErrAPIKeyAlreadyConfirmed")
	})
}
