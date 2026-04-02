package users

import (
	"context"
	"github.com/katharinasick/clubrizer/internal/apperrors"
)

// SignInOrRegisterRequest and SignInOrRegisterResponse are kept temporarily so the
// route wiring in server.go compiles. They will be removed in Task 6 when the
// Google OAuth endpoint is deleted and replaced by the OTP endpoints.
type SignInOrRegisterRequest struct {
	IdToken string `json:"idToken" validate:"required"`
}

type SignInOrRegisterResponse struct {
	AccessToken string `json:"accessToken"`
}

func (s *Service) SignInOrRegister(_ context.Context, _ SignInOrRegisterRequest) (*SignInOrRegisterResponse, *RefreshTokenInfo, error) {
	return nil, nil, apperrors.NewBadRequest("Google sign-in has been removed. Use /auth/otp instead.", nil)
}
