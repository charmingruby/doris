package dynamo

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/charmingruby/doris/lib/instrumentation"
)

type Client struct {
	Client *dynamodb.Client
	logger *instrumentation.Logger
}

type ConnectionInput struct {
	Region string
}

func New(logger *instrumentation.Logger, in ConnectionInput) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration, %v", err)
	}

	cfg.Region = in.Region

	client := dynamodb.NewFromConfig(cfg)

	return &Client{
		Client: client,
		logger: logger,
	}, nil
}
