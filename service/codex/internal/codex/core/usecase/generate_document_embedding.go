package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/model"
	"github.com/charmingruby/doris/service/codex/internal/codex/core/repository"
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

	chunks, err := u.chunkText(rawContent)
	if err != nil {
		return err
	}

	var embeddings [][]float64
	for _, chunk := range chunks {
		embedding, err := u.generateEmbeddingFromChunk(ctx, chunk)
		if err != nil {
			return err
		}

		embeddings = append(embeddings, embedding)
	}

	if err := u.txManager.Transact(func(txManager repository.TransactionManager) error {
		for _, embedding := range embeddings {
			chunk := model.NewCodexDocumentChunk(model.CodexDocumentChunkInput{
				CodexDocumentID: in.DocumentID,
				Embedding:       embedding,
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

func (u *UseCase) generateEmbeddingFromChunk(ctx context.Context, text string) ([]float64, error) {
	reqBody := OllamaEmbeddingRequest{
		Model:  "nomic-embed-text",
		Prompt: text,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, custom_err.NewErrDatasourceOperationFailed("marshal embedding request", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:11434/api/embeddings", strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, custom_err.NewErrDatasourceOperationFailed("create embedding request", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, custom_err.NewErrDatasourceOperationFailed("send embedding request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, custom_err.NewErrDatasourceOperationFailed("embedding request failed", fmt.Errorf("status code: %d", resp.StatusCode))
	}

	var embeddingResp OllamaEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embeddingResp); err != nil {
		return nil, custom_err.NewErrDatasourceOperationFailed("decode embedding response", err)
	}

	return embeddingResp.Embedding, nil
}
