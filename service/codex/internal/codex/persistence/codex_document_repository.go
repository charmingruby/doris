package persistence

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/persistence/postgres"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/model"
	"github.com/jmoiron/sqlx"
)

const (
	createCodexDocument = "create codex document"
)

func codexDocumentQueries() map[string]string {
	return map[string]string{
		createCodexDocument: `INSERT INTO codex_documents (id, codex_id, title, image_url, status, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
	}
}

type CodexDocumentRepository struct {
	db    postgres.Database
	stmts map[string]*sqlx.Stmt
}

func NewCodexDocumentRepository(db postgres.Database) (*CodexDocumentRepository, error) {
	stmts := make(map[string]*sqlx.Stmt)

	for queryName, statement := range codexDocumentQueries() {
		stmt, err := db.Preparex(statement)
		if err != nil {
			return nil, postgres.NewPreparationErr(queryName, "codex document", err)
		}

		stmts[queryName] = stmt
	}

	return &CodexDocumentRepository{
		db:    db,
		stmts: stmts,
	}, nil
}

func (r *CodexDocumentRepository) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]

	if !ok {
		return nil, postgres.NewStatementNotPreparedErr(queryName, "codex document")
	}

	return stmt, nil
}

func (r *CodexDocumentRepository) Create(ctx context.Context, codexDocument model.CodexDocument) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(createCodexDocument)
	if err != nil {
		return err
	}

	if _, err := stmt.ExecContext(ctx,
		codexDocument.ID,
		codexDocument.CodexID,
		codexDocument.Title,
		codexDocument.ImageURL,
		codexDocument.Status,
		codexDocument.CreatedAt,
	); err != nil {
		return err
	}

	return nil
}
