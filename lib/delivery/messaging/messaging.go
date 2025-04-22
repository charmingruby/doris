package messaging

import "context"

type Publisher interface {
	Publish(ctx context.Context, topic string, message []byte) error
	Close(ctx context.Context) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, topic string, handler func(message []byte) error) error
	Close(ctx context.Context) error
}
