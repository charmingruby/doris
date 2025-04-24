package nats

import (
	"context"
	"errors"

	"github.com/charmingruby/doris/lib/instrumentation/logger"
	"github.com/nats-io/nats.go"
)

type Publisher struct {
	logger    *logger.Logger
	nc        *nats.Conn
	js        nats.JetStreamContext
	stream    string
	brokerURL string
}

func NewPublisher(opts ...PublisherOpt) (*Publisher, error) {
	pub := &Publisher{}

	for _, opt := range opts {
		opt(pub)
	}

	if err := pub.validate(); err != nil {
		return nil, err
	}

	nc, err := nats.Connect(pub.brokerURL)
	if err != nil {
		return nil, err
	}

	pub.nc = nc

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	pub.js = js

	if _, err := prepareStream(pub.js, pub.stream); err != nil {
		return nil, err
	}

	return pub, nil
}

func (p *Publisher) Publish(ctx context.Context, topic string, message []byte) error {
	if _, err := prepareSubject(p.js, p.stream, topic); err != nil {
		return err
	}

	if err := p.nc.Publish(topic, message); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		p.logger.Error("message not published", "topic", topic, "error", ctx.Err())
		return ctx.Err()
	default:
		p.logger.Debug("message published", "topic", topic)
		return nil
	}
}

func (s *Publisher) Close(ctx context.Context) error {
	s.nc.Close()

	return nil
}

type PublisherOpt func(*Publisher)

func WithLogger(logger *logger.Logger) PublisherOpt {
	return func(p *Publisher) {
		p.logger = logger
	}
}

func WithBrokerURL(brokerURL string) PublisherOpt {
	return func(p *Publisher) {
		p.brokerURL = brokerURL
	}
}

func WithStream(stream string) PublisherOpt {
	return func(p *Publisher) {
		p.stream = stream
	}
}

func (p *Publisher) validate() error {
	if p.logger == nil {
		return errors.New("logger is required")
	}

	if p.brokerURL == "" {
		p.brokerURL = DEFAULT_BROKER_URL
	}

	if p.stream == "" {
		return errors.New("stream is required")
	}

	return nil
}
