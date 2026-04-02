package users

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/katharinasick/clubrizer/internal/apperrors"
)

type RequestOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (s *Service) RequestOTP(ctx context.Context, req RequestOTPRequest) error {
	count, err := s.store.countRecentOTPRequests(ctx, req.Email)
	if err != nil {
		return err
	}
	if count >= 3 {
		return apperrors.NewTooManyRequests("please wait before requesting a new code")
	}

	if err := s.store.invalidateActiveOTPs(ctx, req.Email); err != nil {
		return err
	}

	code, err := generateOTPCode()
	if err != nil {
		return err
	}

	hash := sha256.Sum256([]byte(code))
	codeHash := hex.EncodeToString(hash[:])

	if err := s.store.createOTPToken(ctx, req.Email, codeHash, time.Now().Add(time.Hour)); err != nil {
		return err
	}

	return s.email.SendOTP(ctx, req.Email, code)
}

func generateOTPCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP code: %w", err)
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
