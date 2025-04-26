package service

import (
	"context"
)

type NotifyApiKeyActivationInput struct{}

func (s *Service) NotifyApiKeyActivation(ctx context.Context, in NotifyApiKeyActivationInput) error {
	return nil
}
