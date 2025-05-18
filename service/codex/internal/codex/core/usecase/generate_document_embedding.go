package usecase

import "context"

type GenerateDocumentEmbeddingInput struct {
	DocumentID string
}

type GenerateDocumentEmbeddingOutput struct {
	Embedding []float64
}

func (u *UseCase) GenerateDocumentEmbedding(ctx context.Context, in GenerateDocumentEmbeddingInput) (GenerateDocumentEmbeddingOutput, error) {
	return GenerateDocumentEmbeddingOutput{}, nil
}
