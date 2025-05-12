package persistence

import (
	"github.com/charmingruby/doris/lib/persistence/postgres"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/repository"
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
		quotaRepo, err := NewQuotaRepository(tx)
		if err != nil {
			return err
		}

		quotaUsageRepo, err := NewQuotaUsageRepository(tx)
		if err != nil {
			return err
		}

		return txFunc(repository.TransactionManager{
			QuotaRepo:      quotaRepo,
			QuotaUsageRepo: quotaUsageRepo,
		})
	})
}
