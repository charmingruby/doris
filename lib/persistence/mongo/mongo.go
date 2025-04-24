package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	db *mongo.Database
}

func New(url, database string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clOpts := options.Client().ApplyURI(url)

	cl, err := mongo.Connect(ctx, clOpts)
	if err != nil {
		return nil, err
	}

	if err := cl.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := cl.Database(database)

	return &Client{db: db}, nil
}

func (c *Client) Close(ctx context.Context) error {
	return c.db.Client().Disconnect(ctx)
}
