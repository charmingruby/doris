package usecase

import "context"

type ConfirmAPIKeyInput struct{}

func (u *UseCase) ConfirmAPIKey(ctx context.Context, in ConfirmAPIKeyInput) error {
	return nil
}
