package users

import (
	"context"
	"encoding/json"
	"github.com/katharinasick/clubrizer/internal/apperrors"
	"google.golang.org/api/idtoken"
)

type SignInOrRegisterRequest struct {
	IdToken string `json:"idToken" validate:"required"`
}

type SignInOrRegisterResponse struct {
	AccessToken string `json:"accessToken"`
}

func (s *Service) SignInOrRegister(ctx context.Context, r SignInOrRegisterRequest) (*SignInOrRegisterResponse, *RefreshTokenInfo, error) {
	// Validate & parse Google ID Token
	token, err := idtoken.Validate(ctx, r.IdToken, s.cfg.OAuth.Google.ClientId)
	if err != nil {
		return nil, nil, apperrors.NewUnauthorized("invalid id token", err)
	}

	c, err := s.parseClaims(token.Claims)
	if err != nil {
		return nil, nil, apperrors.NewBadRequest("failed to parse token googleTokenClaims", err)
	}

	// Create user if it doesn't exist yet
	isNewUser := false
	u, err := s.store.getUserByMail(ctx, c.Email)
	if err != nil {
		if apperrors.Is(err, apperrors.NotFound) {
			isNewUser = true
			u, err = s.store.createUser(ctx, c)
			if err != nil {
				return nil, nil, err
			}
		} else {
			return nil, nil, err
		}
	}

	// Return refresh token as cookie and the access token in the response body
	accessToken, refreshToken, refreshTokenExpiresAt, err := s.generateTokens(u, isNewUser)
	return &SignInOrRegisterResponse{accessToken},
		&RefreshTokenInfo{
			Token:   refreshToken,
			Expires: refreshTokenExpiresAt,
		},
		nil
}

func (s *Service) parseClaims(c map[string]interface{}) (*googleTokenClaims, error) {
	jsonString, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	parsedClaims := googleTokenClaims{}
	err = json.Unmarshal(jsonString, &parsedClaims)
	if err != nil {
		return nil, err
	}
	return &parsedClaims, nil
}
