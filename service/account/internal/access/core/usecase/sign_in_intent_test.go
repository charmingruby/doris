package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/id"
	"github.com/charmingruby/doris/service/account/internal/access/core/model"
)

func (s *Suite) Test_SignInIntent() {
	validInput := SignInIntentInput{
		Email: "john.doe@example.com",
	}

	expirationDelay := 30 * time.Minute

	dummyAPIKey := *model.NewAPIKey(model.APIKeyInput{
		FirstName: "John",
		LastName:  "Doe",
		Email:     validInput.Email,
		Key:       id.New(),
	})

	s.Run("it should dispatch an otp for sign in", func() {
		validAPIKey := dummyAPIKey
		validAPIKey.Status = model.API_KEY_STATUS_ACTIVE

		err := s.apiKeyRepo.Create(context.Background(), validAPIKey)
		s.NoError(err)

		err = s.uc.SignInIntent(context.Background(), validInput)
		s.NoError(err)
		s.Len(s.otpRepo.Items, 1)
		s.Len(s.evtHandler.Pub.Messages, 1)

		otp := s.otpRepo.Items[0]
		expectedExpiration := time.Now().Add(expirationDelay)

		timeDiff := otp.ExpiresAt.Sub(expectedExpiration)

		s.True(timeDiff < time.Second && timeDiff > -time.Second, "expiration time should be within 1 second of expected time")
	})

	s.Run("it should be not able to dispatch an otp for sign in if there is an error inside the transaction", func() {
		validAPIKey := dummyAPIKey
		validAPIKey.Status = model.API_KEY_STATUS_ACTIVE

		err := s.apiKeyRepo.Create(context.Background(), validAPIKey)
		s.NoError(err)

		s.evtHandler.Pub.IsHealthy = false

		err = s.uc.SignInIntent(context.Background(), validInput)

		s.Error(err)
		s.Equal(0, len(s.evtHandler.Pub.Messages))
		s.Equal(0, len(s.otpRepo.Items))
	})

	s.Run("it should be not able to dispatch an otp for sign in if the api key does not exists", func() {
		err := s.uc.SignInIntent(context.Background(), validInput)
		s.Error(err)

		var errResourceNotFound *custom_err.ErrResourceNotFound
		s.True(errors.As(err, &errResourceNotFound), "error should be of type ErrResourceNotFound")
	})

	s.Run("it should be not able to dispatch an otp for sign in if datasource fails", func() {
		validAPIKey := dummyAPIKey
		validAPIKey.Status = model.API_KEY_STATUS_ACTIVE

		err := s.apiKeyRepo.Create(context.Background(), validAPIKey)
		s.NoError(err)

		s.otpRepo.IsHealthy = false

		err = s.uc.SignInIntent(context.Background(), validInput)
		s.Error(err)

		var dsErr *custom_err.ErrDatasourceOperationFailed
		s.True(errors.As(err, &dsErr), "error should be of type ErrDatasourceOperationFailed")
	})

	s.Run("it should be not able to dispatch an otp for sign in if the messaging fails", func() {
		validAPIKey := dummyAPIKey
		validAPIKey.Status = model.API_KEY_STATUS_ACTIVE

		err := s.apiKeyRepo.Create(context.Background(), validAPIKey)
		s.NoError(err)

		s.evtHandler.Pub.IsHealthy = false

		err = s.uc.SignInIntent(context.Background(), validInput)
		s.Error(err)
		s.Equal(0, len(s.evtHandler.Pub.Messages))

		var errMessaging *custom_err.ErrMessagingWrapper
		s.True(errors.As(err, &errMessaging), "error should be of type ErrMessagingWrapper")
	})

	s.Run("it should be not able to dispatch an otp for sign in if api key has insufficient permission", func() {
		invalidAPIKey := dummyAPIKey
		invalidAPIKey.Status = model.API_KEY_STATUS_PENDING

		err := s.apiKeyRepo.Create(context.Background(), invalidAPIKey)
		s.NoError(err)

		err = s.uc.SignInIntent(context.Background(), validInput)
		s.Error(err)

		var errInsufficientPermission *custom_err.ErrInsufficientPermission
		s.True(errors.As(err, &errInsufficientPermission), "error should be of type ErrInsufficientPermission")
	})
}
