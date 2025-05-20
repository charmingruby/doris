package persistence

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/persistence/postgres"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/model"
	"github.com/jmoiron/sqlx"
)

const (
	createQA = "create qa"
)

func qaQueries() map[string]string {
	return map[string]string{
		createQA: `INSERT INTO qas (id, codex_id, question, answer, created_at) VALUES ($1, $2, $3, $4, $5)`,
	}
}

type QARepository struct {
	db    postgres.Database
	stmts map[string]*sqlx.Stmt
}

func NewQARepository(db postgres.Database) (*QARepository, error) {
	stmts := make(map[string]*sqlx.Stmt)

	for queryName, statement := range qaQueries() {
		stmt, err := db.Preparex(statement)
		if err != nil {
			return nil, postgres.NewPreparationErr(queryName, "qa", err)
		}

		stmts[queryName] = stmt
	}

	return &QARepository{
		db:    db,
		stmts: stmts,
	}, nil
}

func (r *QARepository) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]

	if !ok {
		return nil, postgres.NewStatementNotPreparedErr(queryName, "qa")
	}

	return stmt, nil
}

func (r *QARepository) Create(ctx context.Context, qa model.QA) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(createQA)
	if err != nil {
		return err
	}

	if _, err := stmt.ExecContext(ctx,
		qa.ID,
		qa.CodexID,
		qa.Question,
		qa.Answer,
		qa.CreatedAt,
	); err != nil {
		return err
	}

	return nil
}
