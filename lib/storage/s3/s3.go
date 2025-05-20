package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/charmingruby/doris/lib/instrumentation"
)

var (
	ErrUploadFailed   = errors.New("upload failed")
	ErrDownloadFailed = errors.New("download failed")
	ErrInvalidConfig  = errors.New("invalid configuration")
	ErrInvalidFile    = errors.New("invalid file")
)

type Client struct {
	client *s3.Client
	region string
	logger *instrumentation.Logger
}

func New(logger *instrumentation.Logger, region string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		logger.Error("failed to load configuration", "error", err)

		return nil, ErrInvalidConfig
	}

	client := s3.NewFromConfig(cfg)

	return &Client{
		client: client,
		region: region,
		logger: logger,
	}, nil
}

func (c *Client) Upload(ctx context.Context, destination string, key string, file io.Reader) (string, error) {
	c.logger.Debug("uploading file", "destination", destination, "key", key)

	src, err := io.ReadAll(file)
	if err != nil {
		c.logger.Error("failed to read file", "error", err)

		return "", ErrInvalidFile
	}

	if _, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(destination),
		Key:    aws.String(key),
		Body:   bytes.NewReader(src),
	}); err != nil {
		c.logger.Error("failed to upload file", "error", err)

		return "", ErrUploadFailed
	}

	c.logger.Debug("file uploaded", "destination", destination, "key", key)

	return c.bucketFileURL(destination, key), nil
}

func (c *Client) Download(ctx context.Context, source string, url string) (io.Reader, error) {
	c.logger.Debug("downloading file", "source", source, "url", url)

	key, err := c.extractKeyFromURL(url)
	if err != nil {
		c.logger.Error("failed to extract key from URL", "error", err)

		return nil, ErrInvalidFile
	}

	result, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(source),
		Key:    aws.String(key),
	})
	if err != nil {
		c.logger.Error("failed to download file", "error", err)

		return nil, ErrDownloadFailed
	}

	c.logger.Debug("file downloaded", "source", source, "url", url)

	return result.Body, nil
}

func (c *Client) bucketFileURL(destination string, key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", destination, c.region, key)
}

func (c *Client) extractKeyFromURL(fileURL string) (string, error) {
	u, err := url.Parse(fileURL)
	if err != nil {
		c.logger.Error("failed to parse URL", "error", err)

		return "", ErrInvalidFile
	}

	return strings.TrimPrefix(u.Path, "/"), nil
}
