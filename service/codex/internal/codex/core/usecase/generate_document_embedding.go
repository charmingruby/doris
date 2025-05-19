package usecase

import (
	"context"
	"io"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/tmc/langchaingo/textsplitter"
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

type OllamaEmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type OllamaEmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
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

	contentBytes, err := io.ReadAll(doc)
	if err != nil {
		return err
	}

	rawContent := string(contentBytes)

	_, err = u.chunkText(rawContent)
	if err != nil {
		return err
	}

	return nil
}

func (u *UseCase) chunkText(text string) ([]string, error) {
	splitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(1000),
		textsplitter.WithChunkOverlap(200),
	)

	chunks, err := splitter.SplitText(text)
	if err != nil {
		return nil, err
	}

	return chunks, nil
}
