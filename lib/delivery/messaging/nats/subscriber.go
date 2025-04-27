package nats

import (
	"context"

	"github.com/charmingruby/doris/lib/instrumentation"
	"github.com/nats-io/nats.go"
)

type Subscriber struct {
	logger *instrumentation.Logger
	js     nats.JetStreamContext
	nc     *nats.Conn
	cfg    *Config
}

func NewSubscriber(logger *instrumentation.Logger, opts ...ConfigOpt) (*Subscriber, error) {
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

	sub := &Subscriber{
		logger: logger,
		js:     js,
		nc:     nc,
		cfg:    cfg,
	}

	return sub, nil
}

func (s *Subscriber) Subscribe(ctx context.Context, topic string, handler func(message []byte) error) error {
	if _, err := prepareSubject(s.js, s.cfg.stream, topic); err != nil {
		return err
	}

	_, err := s.js.Subscribe(topic, func(msg *nats.Msg) {
		s.logger.Debug("received message", "topic", topic)

		if err := handler(msg.Data); err != nil {
			s.logger.Error("failed to handle message", "error", err)

			if msg.Nak() != nil {
				s.logger.Error("failed to nack message", "error", err)
			}
		}

		if err := msg.Ack(); err != nil {
			s.logger.Error("failed to ack message", "error", err)
		}

		s.logger.Debug("message acknowledged", "topic", topic)
	})

	return err

}

func (s *Subscriber) Close(ctx context.Context) error {
	s.nc.Close()

	select {
	case <-ctx.Done():
		s.logger.Error("failed to close nats subscriber", "error", ctx.Err())
		return ctx.Err()
	default:
		s.logger.Debug("nats subscriber closed")
		return nil
	}
}
