package persistence

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/persistence/postgres"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/model"
	"github.com/jmoiron/sqlx"
)

const (
	createCodexDocumentChunk = "create codex document chunk"
	findSimilarChunks        = "find similar chunks"
)

func codexDocumentChunkQueries() map[string]string {
	return map[string]string{
		createCodexDocumentChunk: `INSERT INTO codex_document_chunks (id, codex_document_id, embedding, content, created_at) VALUES ($1, $2, $3, $4, $5)`,
		findSimilarChunks:        `SELECT cdc.* FROM codex_document_chunks cdc JOIN codex_documents cd ON cdc.codex_document_id = cd.id WHERE cd.codex_id = $1 ORDER BY cdc.embedding <=> $2 LIMIT $3`,
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

	vector := postgres.ParseEmbedding(codexDocumentChunk.Embedding)

	if _, err := stmt.ExecContext(ctx,
		codexDocumentChunk.ID,
		codexDocumentChunk.CodexDocumentID,
		vector,
		codexDocumentChunk.Content,
		codexDocumentChunk.CreatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (r *CodexDocumentChunkRepository) FindSimilarChunks(ctx context.Context, codexID string, embedding []float64, limit int) ([]model.CodexDocumentChunk, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(findSimilarChunks)
	if err != nil {
		return nil, err
	}

	vector := postgres.ParseEmbedding(embedding)

	rows, err := stmt.QueryContext(ctx, codexID, vector, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunks []model.CodexDocumentChunk
	for rows.Next() {
		var chunk model.CodexDocumentChunk
		var embeddingBytes []uint8

		if err := rows.Scan(
			&chunk.ID,
			&chunk.CodexDocumentID,
			&embeddingBytes,
			&chunk.Content,
			&chunk.CreatedAt,
		); err != nil {
			return nil, err
		}

		embedding, err := postgres.ParseEmbeddingFromBytes(embeddingBytes)
		if err != nil {
			return nil, err
		}
		chunk.Embedding = embedding

		chunks = append(chunks, chunk)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chunks, nil
}
