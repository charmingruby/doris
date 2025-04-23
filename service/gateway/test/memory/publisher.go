package memory

import (
	"context"
	"errors"
)

type message struct {
	Content []byte
}

type Publisher struct {
	Messages  []message
	IsHealthy bool
}

func NewPublisher() *Publisher {
	return &Publisher{
		Messages:  []message{},
		IsHealthy: true,
	}
}

func (p *Publisher) Publish(ctx context.Context, topic string, msg []byte) error {
	if !p.IsHealthy {
		return errors.New("publisher is not healthy")
	}

	p.Messages = append(p.Messages, message{Content: msg})

	return nil
}

func (p *Publisher) Close(ctx context.Context) error {
	if !p.IsHealthy {
		return errors.New("publisher is not healthy")
	}

	return nil
}
