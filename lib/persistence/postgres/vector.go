package postgres

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseEmbedding(embedding []float64) string {
	vectorStr := make([]string, len(embedding))
	for i, v := range embedding {
		vectorStr[i] = fmt.Sprintf("%f", v)
	}

	return fmt.Sprintf("[%s]", strings.Join(vectorStr, ","))
}

func ParseEmbeddingFromBytes(bytes []uint8) ([]float64, error) {
	str := string(bytes)
	str = strings.TrimPrefix(str, "[")
	str = strings.TrimSuffix(str, "]")

	values := strings.Split(str, ",")

	embedding := make([]float64, len(values))
	for i, v := range values {
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)

		if err != nil {
			return nil, fmt.Errorf("failed to parse embedding value: %w", err)
		}

		embedding[i] = f
	}

	return embedding, nil
}
