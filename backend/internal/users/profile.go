package users

import (
	"context"
	"io"
)

func (s *Service) GetMyRoles(ctx context.Context) ([]*Role, error) {
	userID := ctx.Value(s.cfg.JWT.User.Key).(*Claims).ID
	roles, err := s.store.getRolesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if roles == nil {
		roles = []*Role{}
	}
	return roles, nil
}

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

	accessToken, refreshToken, refreshTokenExpiresAt, err := s.generateTokens(u)
	if err != nil {
		return nil, nil, err
	}

	return &UpdateProfileResponse{AccessToken: accessToken},
		&RefreshTokenInfo{Token: refreshToken, Expires: refreshTokenExpiresAt},
		nil
}

type UpdateProfilePictureResponse struct {
	AccessToken string `json:"accessToken"`
}

func (s *Service) UpdateProfilePicture(ctx context.Context, contentType string, data io.Reader) (*UpdateProfilePictureResponse, *RefreshTokenInfo, error) {
	claims := ctx.Value(s.cfg.JWT.User.Key).(*Claims)

	url, err := s.storageClient.UploadProfilePicture(ctx, contentType, data)
	if err != nil {
		return nil, nil, err
	}

	u, err := s.store.updateUserPicture(ctx, claims.ID, url)
	if err != nil {
		return nil, nil, err
	}

	accessToken, refreshToken, refreshTokenExpiresAt, err := s.generateTokens(u)
	if err != nil {
		return nil, nil, err
	}

	return &UpdateProfilePictureResponse{AccessToken: accessToken},
		&RefreshTokenInfo{Token: refreshToken, Expires: refreshTokenExpiresAt},
		nil
}
