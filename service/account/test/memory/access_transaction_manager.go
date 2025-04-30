package memory

import (
	"github.com/charmingruby/doris/service/account/internal/access/core/model"
	"github.com/charmingruby/doris/service/account/internal/access/core/repository"
)

type TransactionManager struct {
	apiKeyRepo *APIKeyRepository
	otpRepo    *OTPRepository
}

func NewTransactionManager(apiKeyRepo *APIKeyRepository, otpRepo *OTPRepository) *TransactionManager {
	return &TransactionManager{
		apiKeyRepo: apiKeyRepo,
		otpRepo:    otpRepo,
	}
}

func (t *TransactionManager) Transact(fn func(repos repository.TransactionManager) error) error {
	apiKeySnapshot := make([]model.APIKey, len(t.apiKeyRepo.Items))
	copy(apiKeySnapshot, t.apiKeyRepo.Items)

	otpSnapshot := make([]model.OTP, len(t.otpRepo.Items))
	copy(otpSnapshot, t.otpRepo.Items)

	err := fn(repository.TransactionManager{
		APIKeyRepo: t.apiKeyRepo,
		OTPRepo:    t.otpRepo,
	})

	if err != nil {
		t.apiKeyRepo.Items = apiKeySnapshot
		t.otpRepo.Items = otpSnapshot

		return err
	}

	return nil
}
