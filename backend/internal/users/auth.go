package users

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/katharinasick/clubrizer/internal/apperrors"
	"time"
)

type RefreshTokenInfo struct {
	Token   string
	Expires time.Time
}

type RefreshTokensResponse struct {
	AccessToken string `json:"accessToken"`
}

func (s *Service) RefreshTokens(ctx context.Context, t RefreshTokenInfo) (*RefreshTokensResponse, *RefreshTokenInfo, error) {
	token, err := jwt.Parse(t.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		}

		return []byte(s.cfg.OAuth.RefreshToken.SecretKey), nil
	})

	if err != nil || !token.Valid {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, nil, apperrors.NewUnauthorized("refresh token expired", err)
		}
		return nil, nil, apperrors.NewUnauthorized("invalid refresh token", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, apperrors.NewBadRequest("invalid refresh token claims", nil)
	}
	subject, err := claims.GetSubject()
	if err != nil {
		return nil, nil, apperrors.NewBadRequest("invalid refresh token claims", err)
	}

	user, err := s.store.getUserById(ctx, subject)
	if err != nil {
		return nil, nil, err
	}

	accessToken, refreshToken, refreshTokenExpiresAt, err := s.generateTokens(user, false)
	if err != nil {
		return nil, nil, err
	}

	return &RefreshTokensResponse{accessToken},
		&RefreshTokenInfo{
			Token:   refreshToken,
			Expires: refreshTokenExpiresAt,
		},
		nil
}

func (s *Service) generateTokens(u *User, isNew bool) (string, string, time.Time, error) {
	registeredClaims := jwt.RegisteredClaims{
		Subject:  u.ID.String(),
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}

	accessToken, err := s.generateAccessToken(u, isNew, registeredClaims)
	if err != nil {
		return "", "", time.Time{}, errors.New("failed to generate access token: " + err.Error())
	}
	refreshToken, refreshTokenExpiresAt, err := s.generateRefreshToken(registeredClaims)
	if err != nil {
		return "", "", time.Time{}, errors.New("failed to generate refresh token: " + err.Error())
	}
	return accessToken, refreshToken, refreshTokenExpiresAt, nil
}

func (s *Service) generateAccessToken(u *User, isNew bool, registeredClaims jwt.RegisteredClaims) (string, error) {
	registeredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(15 * time.Minute))
	registeredClaims.Issuer = s.cfg.OAuth.AccessToken.Issuer

	accessClaims := Claims{
		RegisteredClaims: registeredClaims,
		ID:               u.ID,
		IsNew:            isNew,
		Email:            u.Email,
		GivenName:        u.GivenName,
		FamilyName:       u.FamilyName,
		NickName:         u.NickName,
		Picture:          u.Picture,
		Status:           u.Status,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	return token.SignedString([]byte(s.cfg.OAuth.AccessToken.SecretKey))
}

func (s *Service) generateRefreshToken(registeredClaims jwt.RegisteredClaims) (string, time.Time, error) {
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	registeredClaims.ExpiresAt = jwt.NewNumericDate(expiresAt)
	registeredClaims.Issuer = s.cfg.OAuth.RefreshToken.Issuer

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, registeredClaims)
	signedToken, err := token.SignedString([]byte(s.cfg.OAuth.RefreshToken.SecretKey))
	return signedToken, expiresAt, err
}
