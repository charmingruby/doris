package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/charmingruby/doris/lib/instrumentation"
)

type Client struct {
	client *s3.Client
	logger *instrumentation.Logger
}

func New(logger *instrumentation.Logger, region string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration, %v", err)
	}

	client := s3.NewFromConfig(cfg)

	return &Client{
		client: client,
		logger: logger,
	}, nil
}

func (c *Client) Upload(ctx context.Context, destination string, key string, file io.Reader) error {
	src, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file, %v", err)
	}

	_, err = c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(destination),
		Key:    aws.String(key),
		Body:   bytes.NewReader(src),
	})

	return err
}
