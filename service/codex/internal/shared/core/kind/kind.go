package kind

import "errors"

const (
	QUOTA_LIMIT_DOCUMENT = "DOCUMENT"
	QUOTA_LIMIT_PROMPT   = "PROMPT"
)

var (
	ErrInvalidKind = errors.New("invalid kind")

	validKinds = map[string]struct{}{
		QUOTA_LIMIT_DOCUMENT: {},
		QUOTA_LIMIT_PROMPT:   {},
	}
)

func IsValid(kind string) error {
	_, ok := validKinds[kind]

	if !ok {
		return ErrInvalidKind
	}

	return nil
}
