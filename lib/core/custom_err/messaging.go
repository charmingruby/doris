package custom_err

import "fmt"

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
