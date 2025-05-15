package persistence

import (
	"context"
	"database/sql"
	"time"

	"github.com/charmingruby/doris/lib/persistence/postgres"
	"github.com/charmingruby/doris/service/codex/internal/quota/core/model"
	"github.com/jmoiron/sqlx"
)

const (
	findQuotaSnapshotByCorrelationIDAndKind = "find quota snapshot by correlation id and kind"
	updateCurrentUsage                      = "update current usage"
)

func quotaSnapshotQueries() map[string]string {
	return map[string]string{
		findQuotaSnapshotByCorrelationIDAndKind: `
			SELECT 
				q.correlation_id,
				q.kind,
				q.current_usage,
				q.max_value,
				q.unit,
				q.tier,
				q.status,
				qu.last_reset_at
			FROM quota_usages qu
			JOIN quotas q ON q.id = qu.quota_id
			WHERE qu.correlation_id = $1 AND q.kind = $2 AND qu.is_active = true`,
		updateCurrentUsage: `
			UPDATE quota_usages qu
			SET current_usage = $1 
			JOIN quotas q ON q.id = qu.quota_id
			WHERE qu.correlation_id = $2 AND q.kind = $3 AND qu.is_active = true`,
	}
}

type QuotaSnapshotRepository struct {
	db    postgres.Database
	stmts map[string]*sqlx.Stmt
}

func NewQuotaSnapshotRepository(db postgres.Database) (*QuotaSnapshotRepository, error) {
	stmts := make(map[string]*sqlx.Stmt)

	for queryName, statement := range quotaSnapshotQueries() {
		stmt, err := db.Preparex(statement)
		if err != nil {
			return nil,
				postgres.NewPreparationErr(queryName, "quota snapshot", err)
		}

		stmts[queryName] = stmt
	}

	return &QuotaSnapshotRepository{
		db:    db,
		stmts: stmts,
	}, nil
}

func (r *QuotaSnapshotRepository) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]

	if !ok {
		return nil,
			postgres.NewStatementNotPreparedErr(queryName, "quota snapshot")
	}

	return stmt, nil
}

func (r *QuotaSnapshotRepository) FindByCorrelationIDAndKind(ctx context.Context, correlationID, kind string) (model.QuotaSnapshot, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(findQuotaSnapshotByCorrelationIDAndKind)
	if err != nil {
		return model.QuotaSnapshot{}, err
	}

	var snapshot model.QuotaSnapshot
	var lastResetAt sql.NullTime

	if err := stmt.QueryRowContext(ctx, correlationID, kind).Scan(
		&snapshot.CorrelationID,
		&snapshot.Kind,
		&snapshot.CurrentUsage,
		&snapshot.MaxValue,
		&snapshot.Unit,
		&snapshot.Tier,
		&snapshot.Status,
		&lastResetAt,
	); err != nil {
		return model.QuotaSnapshot{}, err
	}

	if lastResetAt.Valid {
		snapshot.LastResetAt = &lastResetAt.Time
	}

	return snapshot, nil
}

func (r *QuotaSnapshotRepository) UpdateCurrentUsage(ctx context.Context, correlationID, kind string, usage int) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(updateCurrentUsage)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx,
		usage,
		correlationID,
		kind,
	)

	return err
}
