package memory

import (
	"errors"

	"github.com/charmingruby/doris/lib/security"
)

type TokenPayload struct {
	Sub     string
	Payload security.Payload
}

type TokenClient struct {
	Items map[string]TokenPayload
}

func NewTokenClient() *TokenClient {
	return &TokenClient{
		Items: make(map[string]TokenPayload),
	}
}

func (m *TokenClient) Generate(sub string, payload security.Payload) (string, error) {
	token := sub + "-token"

	m.Items[token] = TokenPayload{
		Sub:     sub,
		Payload: payload,
	}

	return sub, nil
}

func (m *TokenClient) Validate(token string) (string, security.Payload, error) {
	payload, exists := m.Items[token]

	if !exists {
		return "", security.Payload{}, errors.New("token not found")
	}

	sub := token[:len(token)-6]

	return sub, payload.Payload, nil
}
