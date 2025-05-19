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
	findCodexDocumentByID = "find codex document by id"
	createCodexDocument   = "create codex document"
	saveCodexDocument     = "save codex document"
)

func codexDocumentQueries() map[string]string {
	return map[string]string{
		findCodexDocumentByID: `SELECT * FROM codex_documents WHERE id = $1`,
		createCodexDocument:   `INSERT INTO codex_documents (id, codex_id, title, image_url, status, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		saveCodexDocument:     `UPDATE codex_documents SET status = $1, updated_at = $2 WHERE id = $3`,
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

func (r *CodexDocumentRepository) FindByID(ctx context.Context, id string) (model.CodexDocument, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(findCodexDocumentByID)
	if err != nil {
		return model.CodexDocument{}, err
	}

	var codexDocument model.CodexDocument
	if err := stmt.GetContext(ctx, &codexDocument, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.CodexDocument{}, nil
		}

		return model.CodexDocument{}, err
	}

	return codexDocument, nil
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

func (r *CodexDocumentRepository) Save(ctx context.Context, codexDocument model.CodexDocument) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(saveCodexDocument)
	if err != nil {
		return err
	}

	if _, err := stmt.ExecContext(ctx,
		codexDocument.Status,
		codexDocument.UpdatedAt,
		codexDocument.ID,
	); err != nil {
		return err
	}

	return nil
}
