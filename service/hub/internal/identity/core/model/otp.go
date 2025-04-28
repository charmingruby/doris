package model

import (
	"crypto/rand"
	"time"

	"github.com/charmingruby/doris/lib/core/id"
)

const (
	OTP_PURPOSE_API_KEY_ACTIVATION = "api_key_activation"
)

type OTPInput struct {
	Purpose       string
	CorrelationID string
	ExpiresAt     time.Time
}

func NewOTP(in OTPInput) (*OTP, error) {
	otpCode, err := generateOTPCode(6)

	if err != nil {
		return nil, err
	}

	return &OTP{
		ID:            id.New(),
		Code:          otpCode,
		Purpose:       in.Purpose,
		CorrelationID: in.CorrelationID,
		ExpiresAt:     in.ExpiresAt,
		CreatedAt:     time.Now(),
	}, nil
}

type OTP struct {
	ID            string    `json:"id"`
	Code          string    `json:"code"`
	Purpose       string    `json:"purpose"`
	CorrelationID string    `json:"correlation_id"`
	ExpiresAt     time.Time `json:"expires_at"`
	CreatedAt     time.Time `json:"created_at"`
}

func generateOTPCode(length int) (string, error) {
	otpChars := "1234567890"

	buf := make([]byte, length)

	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := range length {
		buf[i] = otpChars[int(buf[i])%otpCharsLength]
	}

	return string(buf), nil
}
