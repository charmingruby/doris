package usecase

import (
	"context"
	"io"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/model"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/repository"
)

type GenerateDocumentEmbeddingInput struct {
	DocumentID    string
	CodexID       string
	CorrelationID string
	ImageURL      string
}

type GenerateDocumentEmbeddingOutput struct {
	Embedding []float64
}

type chunkWithEmbedding struct {
	content   string
	embedding []float64
}

func (u *UseCase) GenerateDocumentEmbedding(ctx context.Context, in GenerateDocumentEmbeddingInput) error {
	codex, err := u.codexRepo.FindByIDAndCorrelationID(ctx, in.CodexID, in.CorrelationID)
	if err != nil {
		return custom_err.NewErrDatasourceOperationFailed("find codex by id", err)
	}

	if codex.ID == "" {
		return custom_err.NewErrResourceNotFound("codex")
	}

	codexDocument, err := u.codexDocumentRepo.FindByID(ctx, in.DocumentID)
	if err != nil {
		return custom_err.NewErrDatasourceOperationFailed("find codex document by id", err)
	}

	if codexDocument.ID == "" {
		return custom_err.NewErrResourceNotFound("codex document")
	}

	codexDocument.Status = model.CodexDocumentStatusProcessing
	now := time.Now()
	codexDocument.UpdatedAt = &now
	if err := u.codexDocumentRepo.Save(ctx, codexDocument); err != nil {
		return err
	}

	doc, err := u.storage.Download(ctx, u.embeddingSourceDocsBucket, in.ImageURL)
	if err != nil {
		return err
	}

	contentBytes, err := io.ReadAll(doc)
	if err != nil {
		return err
	}

	rawContent := string(contentBytes)

	chunks, err := u.llm.ChunkText(rawContent)
	if err != nil {
		return err
	}

	var chunksWithEmbeddings []chunkWithEmbedding

	for _, chunk := range chunks {
		embedding, err := u.llm.GenerateEmbedding(ctx, chunk)
		if err != nil {
			return err
		}

		chunksWithEmbeddings = append(chunksWithEmbeddings, chunkWithEmbedding{
			content:   chunk,
			embedding: embedding,
		})
	}

	if err := u.txManager.Transact(func(txManager repository.TransactionManager) error {
		for _, chunkWithEmbedding := range chunksWithEmbeddings {
			chunk := model.NewCodexDocumentChunk(model.CodexDocumentChunkInput{
				CodexDocumentID: in.DocumentID,
				Embedding:       chunkWithEmbedding.embedding,
				Content:         chunkWithEmbedding.content,
			})

			if err := txManager.CodexDocumentChunkRepository.Create(ctx, *chunk); err != nil {
				return err
			}
		}

		codexDocument.Status = model.CodexDocumentStatusReady
		now := time.Now()
		codexDocument.UpdatedAt = &now
		if err := u.codexDocumentRepo.Save(ctx, codexDocument); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
