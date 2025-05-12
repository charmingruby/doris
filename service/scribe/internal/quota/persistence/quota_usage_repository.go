package persistence

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/persistence/postgres"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
	"github.com/jmoiron/sqlx"
)

const (
	createQuotaUsage = "create quota usage"
)

func quotaUsageQueries() map[string]string {
	return map[string]string{
		createQuotaUsage: `INSERT INTO quota_usages (id, correlation_id, quota_id, current_usage, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
	}
}

type QuotaUsageRepository struct {
	db    postgres.Database
	stmts map[string]*sqlx.Stmt
}

func NewQuotaUsageRepository(db postgres.Database) (*QuotaUsageRepository, error) {
	stmts := make(map[string]*sqlx.Stmt)

	for queryName, statement := range quotaUsageQueries() {
		stmt, err := db.Preparex(statement)
		if err != nil {
			return nil,
				postgres.NewPreparationErr(queryName, "quota usage", err)
		}

		stmts[queryName] = stmt
	}

	return &QuotaUsageRepository{
		db:    db,
		stmts: stmts,
	}, nil
}

func (r *QuotaUsageRepository) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]

	if !ok {
		return nil,
			postgres.NewStatementNotPreparedErr(queryName, "quota usage")
	}

	return stmt, nil
}

func (r *QuotaUsageRepository) Create(ctx context.Context, quotaUsage model.QuotaUsage) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(createQuotaUsage)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx,
		quotaUsage.ID,
		quotaUsage.CorrelationID,
		quotaUsage.QuotaID,
		quotaUsage.CurrentUsage,
		quotaUsage.CreatedAt,
	)

	return err
}
