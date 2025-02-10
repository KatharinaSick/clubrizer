package users

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type RefreshTokensResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (s *Service) RefreshTokens(ctx context.Context) (*RefreshTokensResponse, error) {
	return nil, nil
}

func (s *Service) generateTokens(u *User, isNew bool) (string, string, error) {
	now := time.Now()
	registeredClaims := jwt.RegisteredClaims{
		Subject:  u.ID.String(),
		IssuedAt: jwt.NewNumericDate(now),
		Issuer:   s.cfg.OAuth.Issuer,
	}

	accessToken, err := s.generateAccessToken(u, isNew, registeredClaims)
	if err != nil {

	}
	refreshToken, err := s.generateRefreshToken(registeredClaims)
	if err != nil {

	}
	return accessToken, refreshToken, nil
}

func (s *Service) generateAccessToken(u *User, isNew bool, registeredClaims jwt.RegisteredClaims) (string, error) {
	now := time.Now()
	registeredClaims.ExpiresAt = jwt.NewNumericDate(now.Add(15 * time.Minute))
	accessClaims := accessTokenClaims{
		RegisteredClaims: registeredClaims,
		IsNew:            isNew,
		Email:            u.Email,
		GivenName:        u.GivenName,
		FamilyName:       u.FamilyName,
		NickName:         u.NickName,
		Picture:          u.Picture,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	return token.SignedString(s.cfg.OAuth.SecretKey)
}

func (s *Service) generateRefreshToken(registeredClaims jwt.RegisteredClaims) (string, error) {
	now := time.Now()
	registeredClaims.ExpiresAt = jwt.NewNumericDate(now.Add(7 * 24 * time.Hour))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, registeredClaims)
	return token.SignedString(s.cfg.OAuth.SecretKey)
}
