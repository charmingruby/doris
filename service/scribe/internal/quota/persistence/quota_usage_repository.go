package persistence

import (
	"context"
	"database/sql"
	"time"

	"github.com/charmingruby/doris/lib/persistence/postgres"
	"github.com/charmingruby/doris/service/scribe/internal/quota/core/model"
	"github.com/jmoiron/sqlx"
)

const (
	findQuotaUsageByCorrelationIDAndQuotaID = "find quota usage by correlation id and quota id"
	createQuotaUsage                        = "create quota usage"
	saveQuotaUsage                          = "save quota usage"
)

func quotaUsageQueries() map[string]string {
	return map[string]string{
		createQuotaUsage:                        `INSERT INTO quota_usages (id, correlation_id, quota_id, current_usage, is_active, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		saveQuotaUsage:                          `UPDATE quota_usages SET current_usage = $1, is_active = $2, updated_at = $3 WHERE id = $4`,
		findQuotaUsageByCorrelationIDAndQuotaID: `SELECT * FROM quota_usages WHERE correlation_id = $1 AND quota_id = $2`,
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

func (r *QuotaUsageRepository) FindByCorrelationIDAndQuotaID(ctx context.Context, correlationID, quotaID string) (model.QuotaUsage, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(findQuotaUsageByCorrelationIDAndQuotaID)
	if err != nil {
		return model.QuotaUsage{}, err
	}

	var usage model.QuotaUsage
	var lastResetAt, updatedAt sql.NullTime

	if err := stmt.QueryRowContext(ctx, correlationID, quotaID).Scan(
		&usage.ID,
		&usage.CorrelationID,
		&usage.QuotaID,
		&usage.CurrentUsage,
		&usage.IsActive,
		&usage.CreatedAt,
		&lastResetAt,
		&updatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return model.QuotaUsage{}, nil
		}

		return model.QuotaUsage{}, err
	}

	if lastResetAt.Valid {
		usage.LastResetAt = &lastResetAt.Time
	}

	if updatedAt.Valid {
		usage.UpdatedAt = &updatedAt.Time
	}

	return usage, nil
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
		quotaUsage.IsActive,
		quotaUsage.CreatedAt,
	)

	return err
}

func (r *QuotaUsageRepository) Save(ctx context.Context, quotaUsage model.QuotaUsage) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(saveQuotaUsage)
	if err != nil {
		return err
	}

	now := time.Now()
	quotaUsage.UpdatedAt = &now

	_, err = stmt.ExecContext(ctx,
		quotaUsage.CurrentUsage,
		quotaUsage.IsActive,
		quotaUsage.UpdatedAt,
		quotaUsage.ID,
	)

	return err
}
