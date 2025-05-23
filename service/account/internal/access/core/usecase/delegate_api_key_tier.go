package usecase

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/core/privilege"
	"github.com/charmingruby/doris/service/account/internal/access/core/event"
	"github.com/charmingruby/doris/service/account/internal/access/core/repository"
)

type DelegateAPIKeyTierInput struct {
	ManagerAPIKeyID  string `json:"manager_api_key_id"`
	APIKeyIDToChange string `json:"api_key_id"`
	NewTier          string `json:"new_tier"`
}

func (uc *UseCase) DelegateAPIKeyTier(ctx context.Context, in DelegateAPIKeyTierInput) error {
	apiKey, err := uc.apiKeyRepo.FindByID(ctx, in.APIKeyIDToChange)

	if err != nil {
		return custom_err.NewErrDatasourceOperationFailed("find api key by id", err)
	}

	if apiKey.ID == "" {
		return custom_err.NewErrResourceNotFound("api key")
	}

	if apiKey.Tier == in.NewTier {
		return custom_err.NewErrNothingToChange()
	}

	managerAPIKey, err := uc.apiKeyRepo.FindByID(ctx, in.ManagerAPIKeyID)

	if err != nil {
		return custom_err.NewErrDatasourceOperationFailed("find api key by id", err)
	}

	if managerAPIKey.ID == "" {
		return custom_err.NewErrResourceNotFound("api key")
	}

	isAdmin := managerAPIKey.Tier == privilege.TIER_ADMIN

	if !isAdmin &&
		(apiKey.Tier == privilege.TIER_ADMIN || in.NewTier == string(privilege.TIER_ADMIN)) {
		return custom_err.NewErrInsufficientPermission()
	}

	oldTier := apiKey.Tier
	apiKey.Tier = in.NewTier

	if err := apiKey.Validate(); err != nil {
		return custom_err.NewErrInvalidEntity(err.Error())
	}

	if err := uc.txManager.Transact(func(tx repository.TransactionManager) error {
		if err := tx.APIKeyRepo.Save(ctx, apiKey); err != nil {
			return custom_err.NewErrDatasourceOperationFailed("save api key", err)
		}

		event := event.APIKeyDelegated{
			ID:      apiKey.ID,
			NewTier: apiKey.Tier,
			OldTier: oldTier,
			SentAt:  time.Now(),
		}

		if err := uc.event.DispatchAPIKeyDelegated(ctx, event); err != nil {
			return custom_err.NewErrMessagingWrapper(err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
