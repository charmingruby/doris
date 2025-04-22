package nats

import (
	"fmt"
	"slices"
	"strings"

	"github.com/nats-io/nats.go"
)

const (
	DEFAULT_BROKER_URL = "nats://localhost:4222"
)

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
