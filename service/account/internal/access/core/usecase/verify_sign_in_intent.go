package usecase

import (
	"context"
	"time"

	"github.com/charmingruby/doris/lib/core/custom_err"
	"github.com/charmingruby/doris/lib/security"
	"github.com/charmingruby/doris/service/account/internal/access/core/model"
)

type VerifySignInIntentInput struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func (uc *UseCase) VerifySignInIntent(ctx context.Context, in VerifySignInIntentInput) (string, error) {
	ak, err := uc.apiKeyRepo.FindByEmail(ctx, in.Email)

	if err != nil {
		return "", custom_err.NewErrDatasourceOperationFailed("find api key by email", err)
	}

	if ak.ID == "" {
		return "", custom_err.NewErrResourceNotFound("api key")
	}

	hasSufficientPermission := ak.Status == model.API_KEY_STATUS_ACTIVE

	if !hasSufficientPermission {
		return "", custom_err.NewErrInsufficientPermission()
	}

	otp, err := uc.otpRepo.FindMostRecentByCorrelationID(ctx, ak.ID)

	if err != nil {
		return "", custom_err.NewErrDatasourceOperationFailed("find otp by correlation id", err)
	}

	if otp.ID == "" {
		return "", custom_err.NewErrResourceNotFound("otp")
	}

	if otp.Code != in.OTP {
		return "", custom_err.NewErrInvalidOTPCode("does not match")
	}

	if otp.ExpiresAt.Before(time.Now().UTC()) {
		return "", custom_err.NewErrInvalidOTPCode("expired")
	}

	return uc.tokenClient.Generate(ak.ID, security.Payload{
		Tier: ak.Tier,
	})
}
