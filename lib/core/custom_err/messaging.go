package custom_err

import "fmt"

// ErrMessagingPublishFailed is an error that occurs when a message fails to be published to a messaging system.
type ErrMessagingPublishFailed struct {
	topic       string
	message     []byte
	originalErr error
}

func (e *ErrMessagingPublishFailed) Error() string {
	return fmt.Sprintf("messaging publish failed: `%s`, %s", e.topic, e.originalErr.Error())
}

func NewErrMessagingPublishFailed(topic string, message []byte, originalErr error) error {
	return &ErrMessagingPublishFailed{topic: topic, message: message, originalErr: originalErr}
}

// ErrMessagingWrapper is an error to wrap messaging errors.
type ErrMessagingWrapper struct {
	originalErr error
}

func (e *ErrMessagingWrapper) Error() string {
	return fmt.Sprintf("messaging error: %s", e.originalErr.Error())
}

func NewErrMessagingWrapper(originalErr error) error {
	return &ErrMessagingWrapper{originalErr: originalErr}
}
