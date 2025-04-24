package service

import (
	"context"
	"errors"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/id"
	"github.com/charmingruby/doris/service/gateway/internal/identity/core/model"
)

func (s *Suite) Test_GenerateAPIKey() {
	validInput := GenerateAPIKeyInput{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	expirationDelay := 10 * time.Minute

	dummyAPIKey := *model.NewAPIKey(model.APIKeyInput{
		FirstName:               validInput.FirstName,
		LastName:                validInput.LastName,
		Email:                   validInput.Email,
		Key:                     id.New(),
		ActivationCodeExpiresAt: time.Now().Add(expirationDelay),
	})

	s.Run("it should create a new api key", func() {
		err := s.svc.GenerateAPIKey(context.Background(), validInput)
		s.NoError(err)

		apiKey := s.apiKeyRepo.Items[0]

		s.NotEmpty(apiKey.ID)
		s.Equal(validInput.FirstName, apiKey.FirstName)
		s.Equal(validInput.LastName, apiKey.LastName)
		s.Equal(validInput.Email, apiKey.Email)
		s.Equal(apiKey.Status, model.API_KEY_STATUS_PENDING)

		expectedExpiration := time.Now().Add(expirationDelay)

		timeDiff := apiKey.ActivationCodeExpiresAt.Sub(expectedExpiration)

		s.True(timeDiff < time.Second && timeDiff > -time.Second, "expiration time should be within 1 second of expected time")

		s.Equal(1, len(s.evtHandler.Pub.Messages))
	})

	s.Run("it should return an error if datasource fails", func() {
		s.apiKeyRepo.IsHealthy = false

		err := s.svc.GenerateAPIKey(context.Background(), validInput)
		s.Error(err)

		var dsErr *custom_err.ErrDatasourceOperationFailed
		s.True(errors.As(err, &dsErr), "error should be of type ErrDatasourceOperationFailed")
	})

	s.Run("it should return an error if the api key already exists", func() {
		err := s.apiKeyRepo.Create(context.Background(), dummyAPIKey)
		s.NoError(err)

		err = s.svc.GenerateAPIKey(context.Background(), validInput)
		s.Error(err)

		var errResourceAlreadyExists *custom_err.ErrResourceAlreadyExists
		s.True(errors.As(err, &errResourceAlreadyExists), "error should be of type ErrResourceAlreadyExists")
	})
}
