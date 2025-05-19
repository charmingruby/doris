package persistence

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmingruby/doris/lib/persistence/postgres"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/model"
	"github.com/jmoiron/sqlx"
)

const (
	createCodexDocumentChunk = "create codex document chunk"
)

func codexDocumentChunkQueries() map[string]string {
	return map[string]string{
		createCodexDocumentChunk: `INSERT INTO codex_document_chunks (id, codex_document_id, embedding, created_at) VALUES ($1, $2, $3, $4)`,
	}
}

type CodexDocumentChunkRepository struct {
	db    postgres.Database
	stmts map[string]*sqlx.Stmt
}

func NewCodexDocumentChunkRepository(db postgres.Database) (*CodexDocumentChunkRepository, error) {
	stmts := make(map[string]*sqlx.Stmt)

	for queryName, statement := range codexDocumentChunkQueries() {
		stmt, err := db.Preparex(statement)
		if err != nil {
			return nil, postgres.NewPreparationErr(queryName, "codex document chunk", err)
		}

		stmts[queryName] = stmt
	}

	return &CodexDocumentChunkRepository{
		db:    db,
		stmts: stmts,
	}, nil
}

func (r *CodexDocumentChunkRepository) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]

	if !ok {
		return nil, postgres.NewStatementNotPreparedErr(queryName, "codex document chunk")
	}

	return stmt, nil
}

func (r *CodexDocumentChunkRepository) Create(ctx context.Context, codexDocumentChunk model.CodexDocumentChunk) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(createCodexDocumentChunk)
	if err != nil {
		return err
	}

	vector := parseEmbedding(codexDocumentChunk.Embedding)

	if _, err := stmt.ExecContext(ctx,
		codexDocumentChunk.ID,
		codexDocumentChunk.CodexDocumentID,
		vector,
		codexDocumentChunk.CreatedAt,
	); err != nil {
		return err
	}

	return nil
}

func parseEmbedding(embedding []float64) string {
	vectorStr := make([]string, len(embedding))
	for i, v := range embedding {
		vectorStr[i] = fmt.Sprintf("%f", v)
	}

	return fmt.Sprintf("[%s]", strings.Join(vectorStr, ","))
}
