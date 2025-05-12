package persistence

import (
	"github.com/charmingruby/doris/lib/persistence/postgres"
	"github.com/charmingruby/doris/service/account/internal/access/core/repository"
	"github.com/jmoiron/sqlx"
)

func NewTransactionManager(db *sqlx.DB) (*TransactionManager, error) {
	return &TransactionManager{
		db: db,
	}, nil
}

type TransactionManager struct {
	db *sqlx.DB
}

func (p *TransactionManager) Transact(txFunc func(params repository.TransactionManager) error) error {
	return postgres.RunInTx(p.db, func(tx *sqlx.Tx) error {
		apiKeyRepo, err := NewAPIKeyRepository(tx)
		if err != nil {
			return err
		}

		otpRepo, err := NewOTPRepository(tx)
		if err != nil {
			return err
		}

		return txFunc(repository.TransactionManager{
			APIKeyRepo: apiKeyRepo,
			OTPRepo:    otpRepo,
		})
	})
}
