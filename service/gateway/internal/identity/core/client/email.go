package client

import (
	"context"
)

type Envelope struct {
	From    string
	To      string
	Subject string
	Body    string
}

type EmailClient interface {
	Send(ctx context.Context, en Envelope) error
}
