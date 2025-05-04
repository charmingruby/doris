package service

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/account/internal/access/core/event"
	"github.com/charmingruby/doris/service/account/internal/access/core/model"
	"github.com/charmingruby/doris/service/account/internal/access/core/repository"
)

type DelegateAPIKeyTierInput struct {
	ManagerAPIKeyID  string `json:"manager_api_key_id"`
	APIKeyIDToChange string `json:"api_key_id"`
	NewTier          string `json:"new_tier"`
}

func (s *Service) DelegateAPIKeyTier(ctx context.Context, in DelegateAPIKeyTierInput) error {
	apiKey, err := s.apiKeyRepo.FindByID(ctx, in.APIKeyIDToChange)

	if err != nil {
		return custom_err.NewErrDatasourceOperationFailed("find api key by id", err)
	}

	if apiKey.ID == "" {
		return custom_err.NewErrResourceNotFound("api key")
	}

	if apiKey.Tier == in.NewTier {
		return custom_err.NewErrNothingToChange()
	}

	managerAPIKeyID, err := s.apiKeyRepo.FindByID(ctx, in.ManagerAPIKeyID)

	if err != nil {
		return custom_err.NewErrDatasourceOperationFailed("find api key by id", err)
	}

	if managerAPIKeyID.ID == "" {
		return custom_err.NewErrResourceNotFound("api key")
	}

	isAdmin := managerAPIKeyID.Tier == model.API_KEY_TIER_ADMIN

	if !isAdmin && in.NewTier == string(model.API_KEY_TIER_ADMIN) {
		return custom_err.NewErrInsufficientPermission()
	}

	oldTier := apiKey.Tier
	apiKey.Tier = in.NewTier

	if err := apiKey.Validate(); err != nil {
		return custom_err.NewErrInvalidEntity(err.Error())
	}

	if err := s.txManager.Transact(func(tx repository.TransactionManager) error {
		if err := tx.APIKeyRepo.Update(ctx, apiKey); err != nil {
			return custom_err.NewErrDatasourceOperationFailed("update api key", err)
		}

		event := event.SendNewAPIKeyDelegationMessage{
			ID:      apiKey.ID,
			NewTier: apiKey.Tier,
			OldTier: oldTier,
			SentAt:  time.Now(),
		}

		if err := s.event.SendNewAPIKeyDelegation(ctx, event); err != nil {
			return custom_err.NewErrMessagingWrapper(err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
