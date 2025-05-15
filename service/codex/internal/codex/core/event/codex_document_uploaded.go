package event

import "time"

type CodexDocumentUploaded struct {
	ID            string
	CodexID       string
	CorrelationID string
	ImageURL      string
	SentAt        time.Time
}
