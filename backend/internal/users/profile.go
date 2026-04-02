package users

import (
	"context"
	"io"
)

type UpdateProfileRequest struct {
	FirstName string  `json:"firstName" validate:"required"`
	LastName  string  `json:"lastName" validate:"required"`
	NickName  *string `json:"nickName"`
}

type UpdateProfileResponse struct {
	AccessToken string `json:"accessToken"`
}

func (s *Service) UpdateProfile(ctx context.Context, req UpdateProfileRequest) (*UpdateProfileResponse, *RefreshTokenInfo, error) {
	claims := ctx.Value(s.cfg.JWT.User.Key).(*Claims)

	u, err := s.store.updateUserProfile(ctx, claims.ID, req.FirstName, req.LastName, req.NickName)
	if err != nil {
		return nil, nil, err
	}

	accessToken, refreshToken, refreshTokenExpiresAt, err := s.generateTokens(u, false)
	if err != nil {
		return nil, nil, err
	}

	return &UpdateProfileResponse{AccessToken: accessToken},
		&RefreshTokenInfo{Token: refreshToken, Expires: refreshTokenExpiresAt},
		nil
}

func (s *Service) UpdateProfilePicture(ctx context.Context, contentType string, data io.Reader) error {
	claims := ctx.Value(s.cfg.JWT.User.Key).(*Claims)

	url, err := s.storageClient.UploadProfilePicture(ctx, contentType, data)
	if err != nil {
		return err
	}

	return s.store.updateUserPicture(ctx, claims.ID, url)
}
