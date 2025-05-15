package event

import "context"

type Handler interface {
	DispatchCodexDocumentUploaded(ctx context.Context, message CodexDocumentUploaded) error
}
