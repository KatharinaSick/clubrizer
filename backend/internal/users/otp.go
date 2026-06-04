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

type VerifyOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}

type VerifyOTPResponse struct {
	AccessToken string `json:"accessToken"`
}

func (s *Service) VerifyOTP(ctx context.Context, req VerifyOTPRequest) (*VerifyOTPResponse, *RefreshTokenInfo, error) {
	otp, err := s.store.getActiveOTPByEmail(ctx, req.Email)
	if err != nil {
		if apperrors.Is(err, apperrors.NotFound) {
			return nil, nil, apperrors.NewUnauthorized("invalid or expired code", nil)
		}
		return nil, nil, err
	}

	hash := sha256.Sum256([]byte(req.Code))
	codeHash := hex.EncodeToString(hash[:])

	if codeHash != otp.CodeHash {
		newCount, err := s.store.incrementOTPAttempts(ctx, otp.ID)
		if err != nil {
			return nil, nil, err
		}
		if newCount >= 5 {
			_ = s.store.invalidateOTP(ctx, otp.ID)
		}
		return nil, nil, apperrors.NewUnauthorized("invalid or expired code", nil)
	}

	if err := s.store.invalidateOTP(ctx, otp.ID); err != nil {
		return nil, nil, err
	}

	u, err := s.store.getUserByMail(ctx, req.Email)
	if err != nil {
		if apperrors.Is(err, apperrors.NotFound) {
			u, err = s.store.createUser(ctx, req.Email)
			if err != nil {
				return nil, nil, err
			}
		} else {
			return nil, nil, err
		}
	}

	accessToken, refreshToken, refreshTokenExpiresAt, err := s.generateTokens(u)
	if err != nil {
		return nil, nil, err
	}

	return &VerifyOTPResponse{AccessToken: accessToken},
		&RefreshTokenInfo{Token: refreshToken, Expires: refreshTokenExpiresAt},
		nil
}

func generateOTPCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP code: %w", err)
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
