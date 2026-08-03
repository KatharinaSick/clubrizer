package users

import (
	"context"

	"github.com/katharinasick/clubrizer/internal/apperrors"
)

// SetAccountTypeRequest captures the signup choice. selfParticipates is a pointer so
// the field is required-but-explicit: false ("only my kids") is a valid value and
// must be distinguishable from a missing field.
type SetAccountTypeRequest struct {
	SelfParticipates *bool `json:"selfParticipates" validate:"required"`
}

type SetAccountTypeResponse struct {
	AccessToken string `json:"accessToken"`
}

// SetAccountType records whether the account holder participates themselves. It also
// keeps the member-role grant in sync (guardians are not members) and reissues tokens
// so the updated selfParticipates claim takes effect immediately. Reachable while the
// account is still pending, as part of onboarding.
func (s *Service) SetAccountType(ctx context.Context, req SetAccountTypeRequest) (*SetAccountTypeResponse, *RefreshTokenInfo, error) {
	claims := ctx.Value(s.cfg.JWT.User.Key).(*Claims)

	u, err := s.store.setSelfParticipates(ctx, claims.ID, *req.SelfParticipates)
	if err != nil {
		return nil, nil, err
	}

	accessToken, refreshToken, refreshTokenExpiresAt, err := s.generateTokens(ctx, u)
	if err != nil {
		return nil, nil, err
	}

	return &SetAccountTypeResponse{AccessToken: accessToken},
		&RefreshTokenInfo{Token: refreshToken, Expires: refreshTokenExpiresAt},
		nil
}

type FinishOnboardingResponse struct {
	AccessToken string `json:"accessToken"`
}

// FinishOnboarding moves the account from 'onboarding' to 'pending' — i.e. submits it
// for approval once the account type and kids are set. Only then does an admin see it
// in the approval queue, so they always know what they're approving. A guardian ("only
// my kids") account must have at least one kid; a participating account may have zero.
// New tokens are reissued so the updated status takes effect immediately.
func (s *Service) FinishOnboarding(ctx context.Context) (*FinishOnboardingResponse, *RefreshTokenInfo, error) {
	claims := ctx.Value(s.cfg.JWT.User.Key).(*Claims)

	if claims.Status != "onboarding" {
		return nil, nil, apperrors.NewBadRequest("onboarding is already complete", nil)
	}

	if !claims.SelfParticipates {
		count, err := s.store.countKidsByUserID(ctx, claims.ID)
		if err != nil {
			return nil, nil, err
		}
		if count == 0 {
			return nil, nil, apperrors.NewBadRequest("please add at least one kid before continuing", nil)
		}
	}

	u, transitioned, err := s.store.submitForApproval(ctx, claims.ID)
	if err != nil {
		return nil, nil, err
	}

	// Only notify admins on a real onboarding→pending transition, so a double-submit or a
	// token replay doesn't re-notify. Fire-and-forget — see notifyAdminsOfNewApprovalRequest.
	if transitioned {
		s.notifyAdminsOfNewApprovalRequest(u)
	}

	accessToken, refreshToken, refreshTokenExpiresAt, err := s.generateTokens(ctx, u)
	if err != nil {
		return nil, nil, err
	}

	return &FinishOnboardingResponse{AccessToken: accessToken},
		&RefreshTokenInfo{Token: refreshToken, Expires: refreshTokenExpiresAt},
		nil
}
