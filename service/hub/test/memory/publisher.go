package memory

import (
	"context"
)

type Message struct {
	Content []byte
}

type Publisher struct {
	Messages  []Message
	IsHealthy bool
}

func NewPublisher() *Publisher {
	return &Publisher{
		Messages:  []Message{},
		IsHealthy: true,
	}
}

func (p *Publisher) Publish(ctx context.Context, topic string, msg []byte) error {
	if !p.IsHealthy {
		return ErrUnhealthyDatasource
	}

	p.Messages = append(p.Messages, Message{Content: msg})

	return nil
}

func (p *Publisher) Close(ctx context.Context) error {
	if !p.IsHealthy {
		return ErrUnhealthyDatasource
	}

	return nil
}
