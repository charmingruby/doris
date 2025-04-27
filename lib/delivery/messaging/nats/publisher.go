package nats

import (
	"context"

	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/nats-io/nats.go"
)

type Publisher struct {
	nc     *nats.Conn
	js     nats.JetStreamContext
	logger *instrumentation.Logger
	cfg    *Config
}

func NewPublisher(logger *instrumentation.Logger, opts ...ConfigOpt) (*Publisher, error) {
	cfg := &Config{}

	for _, opt := range opts {
		opt(cfg)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	nc, err := nats.Connect(cfg.brokerURL)
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	if _, err := prepareStream(js, cfg.stream); err != nil {
		return nil, err
	}

	pub := &Publisher{
		nc:     nc,
		js:     js,
		logger: logger,
		cfg:    cfg,
	}

	return pub, nil
}

func (p *Publisher) Publish(ctx context.Context, topic string, message []byte) error {
	if _, err := prepareSubject(p.js, p.cfg.stream, topic); err != nil {
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

func (p *Publisher) Close(ctx context.Context) error {
	p.nc.Close()

	select {
	case <-ctx.Done():
		p.logger.Error("failed to close nats publisher", "error", ctx.Err())
		return ctx.Err()
	default:
		p.logger.Debug("nats publisher closed")
		return nil
	}
}
