package persistence

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/charmingruby/doris/lib/persistence/postgres"
	"github.com/charmingruby/doris/service/codex/internal/quota/core/model"
	"github.com/jmoiron/sqlx"
)

const (
	findQuotaByID          = "find quota by id"
	findManyQuotasByTier   = "find many quotas by tier"
	findQuotaByTierAndKind = "find quota by tier and kind"
	createQuota            = "create quota"
	saveQuota              = "update quota"
)

func quotaQueries() map[string]string {
	return map[string]string{
		findQuotaByID:          `SELECT * FROM quotas WHERE id = $1`,
		findManyQuotasByTier:   `SELECT * FROM quotas WHERE tier = $1`,
		findQuotaByTierAndKind: `SELECT * FROM quotas WHERE tier = $1 AND kind = $2`,
		createQuota:            `INSERT INTO quotas (id, tier, kind, max_value, unit, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		saveQuota:              `UPDATE quotas SET status = $1, tier = $2, kind = $3, max_value = $4, unit = $5, updated_at = $6 WHERE id = $7`,
	}
}

type QuotaRepository struct {
	db    postgres.Database
	stmts map[string]*sqlx.Stmt
}

func NewQuotaRepository(db postgres.Database) (*QuotaRepository, error) {
	stmts := make(map[string]*sqlx.Stmt)

	for queryName, statement := range quotaQueries() {
		stmt, err := db.Preparex(statement)
		if err != nil {
			return nil,
				postgres.NewPreparationErr(queryName, "quota", err)
		}

		stmts[queryName] = stmt
	}

	return &QuotaRepository{
		db:    db,
		stmts: stmts,
	}, nil
}

func (r *QuotaRepository) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]

	if !ok {
		return nil,
			postgres.NewStatementNotPreparedErr(queryName, "quota")
	}

	return stmt, nil
}

func (r *QuotaRepository) FindByID(ctx context.Context, id string) (model.Quota, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(findQuotaByID)
	if err != nil {
		return model.Quota{}, err
	}

	var quota model.Quota
	var updatedAt sql.NullTime

	if err := stmt.QueryRowContext(ctx, id).Scan(
		&quota.ID,
		&quota.Tier,
		&quota.Kind,
		&quota.MaxValue,
		&quota.Unit,
		&quota.Status,
		&quota.CreatedAt,
		&updatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return model.Quota{}, nil
		}

		return model.Quota{}, err
	}

	if updatedAt.Valid {
		quota.UpdatedAt = &updatedAt.Time
	}

	return quota, nil
}

func (r *QuotaRepository) FindByTierAndKind(ctx context.Context, tier, kind string) (model.Quota, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var quota model.Quota
	var updatedAt sql.NullTime

	stmt, err := r.statement(findQuotaByTierAndKind)
	if err != nil {
		return model.Quota{}, err
	}

	if err := stmt.QueryRowContext(ctx, tier, kind).Scan(
		&quota.ID,
		&quota.Tier,
		&quota.Kind,
		&quota.MaxValue,
		&quota.Unit,
		&quota.Status,
		&quota.CreatedAt,
		&updatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return model.Quota{}, nil
		}

		return model.Quota{}, err
	}

	if updatedAt.Valid {
		quota.UpdatedAt = &updatedAt.Time
	}

	return quota, nil
}

func (r *QuotaRepository) FindManyByTier(ctx context.Context, tier string) ([]model.Quota, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(findManyQuotasByTier)
	if err != nil {
		return nil, err
	}

	var quotas []model.Quota

	if err := stmt.SelectContext(ctx, &quotas, tier); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return quotas, nil
}

func (r *QuotaRepository) Create(ctx context.Context, quota model.Quota) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(createQuota)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx,
		quota.ID,
		quota.Tier,
		quota.Kind,
		quota.MaxValue,
		quota.Unit,
		quota.Status,
		quota.CreatedAt,
		quota.UpdatedAt,
	)

	return err
}

func (r *QuotaRepository) Save(ctx context.Context, quota model.Quota) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(saveQuota)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx,
		quota.Status,
		quota.Tier,
		quota.Kind,
		quota.MaxValue,
		quota.Unit,
		quota.UpdatedAt,
		quota.ID,
	)

	return err
}
