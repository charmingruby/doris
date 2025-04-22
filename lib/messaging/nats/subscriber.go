package nats

import (
	"context"

	"github.com/charmingruby/doris/lib/instrumentation/logger"
)

type Subscriber struct {
	logger *logger.Logger
}

func NewSubscriber(logger *logger.Logger) *Subscriber {
	return &Subscriber{
		logger: logger,
	}
}

func (s *Subscriber) Subscribe(ctx context.Context, topic string, handler func(message []byte) error) error {

	return nil

}

func (s *Subscriber) Close(ctx context.Context) error {
	return nil
}
