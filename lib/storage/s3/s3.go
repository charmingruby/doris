package s3

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/charmingruby/doris/lib/instrumentation"
)

type Client struct {
	Client *s3.Client
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
		Client: client,
		logger: logger,
	}, nil
}
