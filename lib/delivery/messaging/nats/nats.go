package nats

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/nats-io/nats.go"
)

const (
	defaultBrokerURL  = "nats://localhost:4222"
	defaultMaxRetries = 3
)

type Config struct {
	stream     string
	brokerURL  string
	maxRetries int
}

type ConfigOpt func(*Config)

func WithBrokerURL(brokerURL string) ConfigOpt {
	return func(p *Config) {
		p.brokerURL = brokerURL
	}
}

func WithMaxRetries(maxRetries int) ConfigOpt {
	return func(p *Config) {
		p.maxRetries = maxRetries
	}
}

func WithStream(stream string) ConfigOpt {
	return func(p *Config) {
		p.stream = stream
	}
}

func (p *Config) validate() error {
	if p.brokerURL == "" {
		p.brokerURL = defaultBrokerURL
	}

	if p.stream == "" {
		return errors.New("stream is required")
	}

	if p.maxRetries <= 0 {
		p.maxRetries = defaultMaxRetries
	}

	return nil
}

func prepareStream(js nats.JetStreamContext, streamName string) (bool, error) {
	if _, err := js.StreamInfo(streamName); err != nil {
		if _, err := js.AddStream(&nats.StreamConfig{
			Name: streamName,
		}); err != nil {
			return false, fmt.Errorf("failed to create stream: %v", err)
		}

		return true, nil
	}

	return false, nil
}

func prepareSubject(js nats.JetStreamContext, streamName string, subject string) (created bool, err error) {
	streamInfo, err := js.StreamInfo(streamName)
	if err != nil {
		return false, fmt.Errorf("stream not found: %v", err)
	}

	if slices.Contains(streamInfo.Config.Subjects, subject) {
		return false, nil
	}

	streamInfo.Config.Subjects = append(streamInfo.Config.Subjects, subject)
	if _, err := js.UpdateStream(&streamInfo.Config); err != nil {
		if strings.Contains(err.Error(), "subjects overlap") {
			return false, nil
		}

		return false, fmt.Errorf("failed to update stream with new subject: %v", err)
	}

	return true, nil
}
