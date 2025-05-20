package client

import "context"

type LLM interface {
	ChunkText(text string) ([]string, error)
	GenerateEmbedding(ctx context.Context, text string) ([]float64, error)
	GenerateCompletion(ctx context.Context, prompt string, limits []string) (string, error)
}
