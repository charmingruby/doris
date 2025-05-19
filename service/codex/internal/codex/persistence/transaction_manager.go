package persistence

import (
	"github.com/charmingruby/doris/lib/persistence/postgres"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/repository"
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
		codexDocumentRepo, err := NewCodexDocumentRepository(tx)
		if err != nil {
			return err
		}

		codexDocumentChunkRepo, err := NewCodexDocumentChunkRepository(tx)
		if err != nil {
			return err
		}

		return txFunc(repository.TransactionManager{
			CodexDocumentRepository:      codexDocumentRepo,
			CodexDocumentChunkRepository: codexDocumentChunkRepo,
		})
	})
}
