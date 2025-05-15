package persistence

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/charmingruby/doris/lib/persistence/postgres"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/model"
	"github.com/jmoiron/sqlx"
)

const (
	findCodexByCorrelationIDAndName = "find codex by correlation id and name"
	createCodex                     = "create codex"
)

func codexQueries() map[string]string {
	return map[string]string{
		findCodexByCorrelationIDAndName: `SELECT * FROM codex WHERE correlation_id = $1 AND name = $2`,
		createCodex:                     `INSERT INTO codex (id, correlation_id, name, description, created_at) VALUES ($1, $2, $3, $4, $5)`,
	}
}

type CodexRepository struct {
	db    postgres.Database
	stmts map[string]*sqlx.Stmt
}

func NewCodexRepository(db postgres.Database) (*CodexRepository, error) {
	stmts := make(map[string]*sqlx.Stmt)

	for queryName, statement := range codexQueries() {
		stmt, err := db.Preparex(statement)
		if err != nil {
			return nil,
				postgres.NewPreparationErr(queryName, "codex", err)
		}

		stmts[queryName] = stmt
	}

	return &CodexRepository{
		db:    db,
		stmts: stmts,
	}, nil
}

func (r *CodexRepository) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]

	if !ok {
		return nil,
			postgres.NewStatementNotPreparedErr(queryName, "codex")
	}

	return stmt, nil
}

func (r *CodexRepository) Create(ctx context.Context, codex model.Codex) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(createCodex)
	if err != nil {
		return err
	}

	if _, err := stmt.ExecContext(ctx, codex.ID, codex.CorrelationID, codex.Name, codex.Description, codex.CreatedAt); err != nil {
		return err
	}

	return nil
}

func (r *CodexRepository) FindByCorrelationIDAndName(ctx context.Context, correlationID, name string) (model.Codex, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(findCodexByCorrelationIDAndName)
	if err != nil {
		return model.Codex{}, err
	}

	var codex model.Codex
	if err := stmt.GetContext(ctx, &codex, correlationID, name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Codex{}, nil
		}

		return model.Codex{}, err
	}

	return codex, nil
}
