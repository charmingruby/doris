package persistence

import (
	"context"
	"database/sql"
	"errors"

	"github.com/charmingruby/doris/lib/persistence/postgres"
	"github.com/charmingruby/doris/service/identity/internal/access/core/model"
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
		createAPIKey:      `INSERT INTO api_keys (id, first_name, last_name, email, key, status) VALUES ($1, $2, $3, $4, $5, $6)`,
		updateAPIKey:      `UPDATE api_keys SET status = $1 WHERE id = $2`,
	}
}

type APIKeyPostgresRepo struct {
	db    *sqlx.DB
	stmts map[string]*sqlx.Stmt
}

func NewAPIKeyPostgresRepo(db *sqlx.DB) (*APIKeyPostgresRepo, error) {
	stmts := make(map[string]*sqlx.Stmt)

	for queryName, statement := range apiKeyQueries() {
		stmt, err := db.Preparex(statement)
		if err != nil {
			return nil,
				postgres.NewPreparationErr(queryName, "api key", err)
		}

		stmts[queryName] = stmt
	}

	return &APIKeyPostgresRepo{
		db:    db,
		stmts: stmts,
	}, nil
}

func (r *APIKeyPostgresRepo) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]

	if !ok {
		return nil,
			postgres.NewStatementNotPreparedErr(queryName, "api key")
	}

	return stmt, nil
}

func (r *APIKeyPostgresRepo) FindByID(ctx context.Context, id string) (model.APIKey, error) {
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

func (r *APIKeyPostgresRepo) FindByEmail(ctx context.Context, email string) (model.APIKey, error) {
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

func (r *APIKeyPostgresRepo) FindByKey(ctx context.Context, key string) (model.APIKey, error) {
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

func (r *APIKeyPostgresRepo) Create(ctx context.Context, apiKey model.APIKey) error {
	stmt, err := r.statement(createAPIKey)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(
		apiKey.ID,
		apiKey.FirstName,
		apiKey.LastName,
		apiKey.Email,
		apiKey.Key,
		apiKey.Status,
	); err != nil {
		return err
	}

	return nil
}

func (r *APIKeyPostgresRepo) Update(ctx context.Context, apiKey model.APIKey) error {
	stmt, err := r.statement(updateAPIKey)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(apiKey.Status, apiKey.ID); err != nil {
		return err
	}

	return nil
}
