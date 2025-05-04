package service

import (
	"context"
	"errors"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/id"
	"github.com/charmingruby/doris/service/account/internal/access/core/model"
)

func (s *Suite) Test_DelegateAPIKeyTier() {
	dummyManagerAPIKey := *model.NewAPIKey(model.APIKeyInput{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe1@example.com",
		Key:       id.New(),
	})
	dummyManagerAPIKey.Status = model.API_KEY_STATUS_ACTIVE
	dummyManagerAPIKey.Tier = model.API_KEY_TIER_ADMIN

	dummyAPIKey := *model.NewAPIKey(model.APIKeyInput{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe2@example.com",
		Key:       id.New(),
	})
	dummyAPIKey.Status = model.API_KEY_STATUS_ACTIVE

	s.Run("it should be able to delegate a new api key tier", func() {
		ctx := context.Background()

		err := s.apiKeyRepo.Create(ctx, dummyManagerAPIKey)
		s.NoError(err)

		err = s.apiKeyRepo.Create(ctx, dummyAPIKey)
		s.NoError(err)

		input := DelegateAPIKeyTierInput{
			ManagerAPIKeyID:  dummyManagerAPIKey.ID,
			APIKeyIDToChange: dummyAPIKey.ID,
			NewTier:          model.API_KEY_TIER_PRO,
		}

		err = s.svc.DelegateAPIKeyTier(ctx, input)
		s.NoError(err)

		modifiedAPIKey, err := s.apiKeyRepo.FindByID(ctx, dummyAPIKey.ID)
		s.NoError(err)

		s.Equal(modifiedAPIKey.Tier, input.NewTier)
	})

	s.Run("it should be not able to delegate a new api key tier if the datasource operation fails", func() {
		ctx := context.Background()

		err := s.apiKeyRepo.Create(ctx, dummyManagerAPIKey)
		s.NoError(err)

		err = s.apiKeyRepo.Create(ctx, dummyAPIKey)
		s.NoError(err)

		s.apiKeyRepo.IsHealthy = false

		input := DelegateAPIKeyTierInput{
			ManagerAPIKeyID:  dummyManagerAPIKey.ID,
			APIKeyIDToChange: dummyAPIKey.ID,
			NewTier:          model.API_KEY_TIER_PRO,
		}

		err = s.svc.DelegateAPIKeyTier(ctx, input)
		s.Error(err)

		var dsErr *custom_err.ErrDatasourceOperationFailed
		s.True(errors.As(err, &dsErr), "error should be of type ErrDatasourceOperationFailed")
	})

	s.Run("it should be not able to delegate a new api key tier if the api key to modify does not exists", func() {
		ctx := context.Background()

		err := s.apiKeyRepo.Create(ctx, dummyManagerAPIKey)
		s.NoError(err)

		input := DelegateAPIKeyTierInput{
			ManagerAPIKeyID:  dummyManagerAPIKey.ID,
			APIKeyIDToChange: "invalid id",
			NewTier:          model.API_KEY_TIER_PRO,
		}

		err = s.svc.DelegateAPIKeyTier(ctx, input)
		s.Error(err)

		var dsErr *custom_err.ErrResourceNotFound
		s.True(errors.As(err, &dsErr), "error should be of type ErrResourceNotFound")
	})

	s.Run("it should be not able to delegate a new api key tier if there is nothing to modify", func() {
		ctx := context.Background()

		err := s.apiKeyRepo.Create(ctx, dummyManagerAPIKey)
		s.NoError(err)

		err = s.apiKeyRepo.Create(ctx, dummyAPIKey)
		s.NoError(err)

		input := DelegateAPIKeyTierInput{
			ManagerAPIKeyID:  dummyManagerAPIKey.ID,
			APIKeyIDToChange: dummyAPIKey.ID,
			NewTier:          dummyAPIKey.Tier,
		}

		err = s.svc.DelegateAPIKeyTier(ctx, input)
		s.Error(err)

		var dsErr *custom_err.ErrNothingToChange
		s.True(errors.As(err, &dsErr), "error should be of type ErrNothingToChange")
	})

	s.Run("it should be not able to delegate a new api key tier if the api key manager does not exists", func() {
		ctx := context.Background()

		err := s.apiKeyRepo.Create(ctx, dummyAPIKey)
		s.NoError(err)

		input := DelegateAPIKeyTierInput{
			ManagerAPIKeyID:  "invalid id",
			APIKeyIDToChange: dummyAPIKey.ID,
			NewTier:          model.API_KEY_TIER_PRO,
		}

		err = s.svc.DelegateAPIKeyTier(ctx, input)
		s.Error(err)

		var dsErr *custom_err.ErrResourceNotFound
		s.True(errors.As(err, &dsErr), "error should be of type ErrResourceNotFound")
	})

	s.Run("it should be not able to delegate a new api key tier if the new tier is admin and manager is not a admin", func() {
		ctx := context.Background()

		managerAPIKey := dummyManagerAPIKey
		managerAPIKey.Tier = model.API_KEY_TIER_MANAGER

		err := s.apiKeyRepo.Create(ctx, managerAPIKey)
		s.NoError(err)

		err = s.apiKeyRepo.Create(ctx, dummyAPIKey)
		s.NoError(err)

		input := DelegateAPIKeyTierInput{
			ManagerAPIKeyID:  managerAPIKey.ID,
			APIKeyIDToChange: dummyAPIKey.ID,
			NewTier:          model.API_KEY_TIER_ADMIN,
		}

		err = s.svc.DelegateAPIKeyTier(ctx, input)
		s.Error(err)

		var dsErr *custom_err.ErrInsufficientPermission
		s.True(errors.As(err, &dsErr), "error should be of type ErrInsufficientPermission")
	})

	s.Run("it should be not able to delegate a new api key tier if the tier does not exists", func() {
		ctx := context.Background()

		err := s.apiKeyRepo.Create(ctx, dummyManagerAPIKey)
		s.NoError(err)

		err = s.apiKeyRepo.Create(ctx, dummyAPIKey)
		s.NoError(err)

		input := DelegateAPIKeyTierInput{
			ManagerAPIKeyID:  dummyManagerAPIKey.ID,
			APIKeyIDToChange: dummyAPIKey.ID,
			NewTier:          "INVALID",
		}

		err = s.svc.DelegateAPIKeyTier(ctx, input)
		s.Error(err)

		var dsErr *custom_err.ErrInvalidEntity
		s.True(errors.As(err, &dsErr), "error should be of type ErrInvalidEntity")
	})

	s.Run("it should be not able to delegate a new api key tier if the messaging fails", func() {
		ctx := context.Background()

		err := s.apiKeyRepo.Create(ctx, dummyManagerAPIKey)
		s.NoError(err)

		err = s.apiKeyRepo.Create(ctx, dummyAPIKey)
		s.NoError(err)

		input := DelegateAPIKeyTierInput{
			ManagerAPIKeyID:  dummyManagerAPIKey.ID,
			APIKeyIDToChange: dummyAPIKey.ID,
			NewTier:          model.API_KEY_TIER_PRO,
		}

		s.evtHandler.Pub.IsHealthy = false

		err = s.svc.DelegateAPIKeyTier(ctx, input)
		s.Error(err)

		var dsErr *custom_err.ErrMessagingWrapper
		s.True(errors.As(err, &dsErr), "error should be of type ErrMessagingWrapper")
	})
}
