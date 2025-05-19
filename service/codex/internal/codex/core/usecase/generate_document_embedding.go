package usecase

import (
	"context"
	"io"

	"github.com/charmingruby/doris/lib/core/custom_err"
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

	doc, err := u.storage.Download(ctx, u.embeddingSourceDocsBucket, in.ImageURL)
	if err != nil {
		return err
	}

	docBytes, err := io.ReadAll(doc)
	if err != nil {
		return custom_err.NewErrDatasourceOperationFailed("read image", err)
	}

	docContent := string(docBytes)

	println(docContent)

	return nil
}
