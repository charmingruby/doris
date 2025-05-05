package persistence

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/charmingruby/doris/lib/persistence/postgres"
	"github.com/charmingruby/doris/service/account/internal/access/core/model"
	"github.com/jmoiron/sqlx"
)

const (
	findAPIKeyByID    = "find api key by id"
	findAPIKeyByKey   = "find api key by key"
	findAPIKeyByEmail = "find api key by email"
	createAPIKey      = "create api key"
	updateAPIKey      = "update api key"
)

func apiKeyQueries() map[string]string {
	return map[string]string{
		findAPIKeyByID:    `SELECT * FROM api_keys WHERE id = $1`,
		findAPIKeyByKey:   `SELECT * FROM api_keys WHERE key = $1`,
		findAPIKeyByEmail: `SELECT * FROM api_keys WHERE email = $1`,
		createAPIKey:      `INSERT INTO api_keys (id, first_name, last_name, email, key, status, tier) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		updateAPIKey:      `UPDATE api_keys SET status = $1, tier = $2 WHERE id = $3`,
	}
}

type APIKeyRepo struct {
	db    postgres.Database
	stmts map[string]*sqlx.Stmt
}

func NewAPIKeyRepo(db postgres.Database) (*APIKeyRepo, error) {
	stmts := make(map[string]*sqlx.Stmt)

	for queryName, statement := range apiKeyQueries() {
		stmt, err := db.Preparex(statement)
		if err != nil {
			return nil,
				postgres.NewPreparationErr(queryName, "api key", err)
		}

		stmts[queryName] = stmt
	}

	return &APIKeyRepo{
		db:    db,
		stmts: stmts,
	}, nil
}

func (r *APIKeyRepo) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]

	if !ok {
		return nil,
			postgres.NewStatementNotPreparedErr(queryName, "api key")
	}

	return stmt, nil
}

func (r *APIKeyRepo) FindByID(ctx context.Context, id string) (model.APIKey, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(findAPIKeyByID)
	if err != nil {
		return model.APIKey{}, err
	}

	var apiKey model.APIKey
	if err := stmt.GetContext(ctx, &apiKey, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.APIKey{}, nil
		}

		return model.APIKey{}, err
	}

	return apiKey, nil
}

func (r *APIKeyRepo) FindByEmail(ctx context.Context, email string) (model.APIKey, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(findAPIKeyByEmail)
	if err != nil {
		return model.APIKey{}, err
	}

	var apiKey model.APIKey
	if err := stmt.GetContext(ctx, &apiKey, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.APIKey{}, nil
		}

		return model.APIKey{}, err
	}

	return apiKey, nil
}

func (r *APIKeyRepo) FindByKey(ctx context.Context, key string) (model.APIKey, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(findAPIKeyByID)
	if err != nil {
		return model.APIKey{}, err
	}

	var apiKey model.APIKey
	if err := stmt.GetContext(ctx, &apiKey, key); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.APIKey{}, nil
		}

		return model.APIKey{}, err
	}

	return apiKey, nil
}

func (r *APIKeyRepo) Create(ctx context.Context, apiKey model.APIKey) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(createAPIKey)
	if err != nil {
		return err
	}

	if _, err := stmt.ExecContext(
		ctx,
		apiKey.ID,
		apiKey.FirstName,
		apiKey.LastName,
		apiKey.Email,
		apiKey.Key,
		apiKey.Status,
		apiKey.Tier,
	); err != nil {
		return err
	}

	return nil
}

func (r *APIKeyRepo) Update(ctx context.Context, apiKey model.APIKey) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(updateAPIKey)
	if err != nil {
		return err
	}

	if _, err := stmt.ExecContext(ctx, apiKey.Status, apiKey.Tier, apiKey.ID); err != nil {
		return err
	}

	return nil
}
