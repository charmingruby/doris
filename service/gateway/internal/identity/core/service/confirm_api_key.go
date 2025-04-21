package service

import "context"

type ConfirmAPIKeyInput struct {
	RequestedAPIKeyID string
	ConfirmationCode  string
}

func (s *Service) ConfirmAPIKey(ctx context.Context, in ConfirmAPIKeyInput) error {
	apiKey, err := s.apiKeyRepo.FindByID(ctx, in.RequestedAPIKeyID)

	return nil
}
