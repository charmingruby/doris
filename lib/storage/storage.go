package storage

import (
	"context"
	"io"
)

type Storage interface {
	Upload(ctx context.Context, destination string, key string, file io.Reader) (string, error)
	Download(ctx context.Context, source string, key string) (io.Reader, error)
}
