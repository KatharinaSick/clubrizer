package users

import (
	"context"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/katharinasick/clubrizer/internal/apperrors"
)

// KidRequest is the payload for adding or editing a kid after approval, where both names
// are required — the same complete-details bar the account holder meets. (The name-only,
// last-name-later step is onboarding, which uses ReplaceKidsRequest instead; family_name
// stays nullable in the schema for those onboarding-created rows.) The approval status is
// managed server-side and the picture has its own endpoint.
type KidRequest struct {
	FirstName string `json:"firstName" validate:"required,notblank"`
	LastName  string `json:"lastName" validate:"required,notblank"`
}

func (s *Service) ListKids(ctx context.Context) ([]*Kid, error) {
	userID := ctx.Value(s.cfg.JWT.User.Key).(*Claims).ID
	kids, err := s.store.getKidsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if kids == nil {
		kids = []*Kid{}
	}
	return kids, nil
}

func (s *Service) AddKid(ctx context.Context, req KidRequest) (*Kid, error) {
	userID := ctx.Value(s.cfg.JWT.User.Key).(*Claims).ID
	return s.store.createKid(ctx, userID, req.FirstName, req.LastName)
}

// ReplaceKidsRequest is the onboarding "add your kids" payload: the whole list of first
// names at once. Blank entries are ignored so trailing empty fields don't create kids.
type ReplaceKidsRequest struct {
	FirstNames []string `json:"firstNames"`
}

// ReplaceKids sets the account's kids to exactly the supplied first names — the
// onboarding screen sends its full list in one request, and this makes the stored kids
// match it (adds new, drops removed) in a single transaction. Gated to onboarding at the
// route level, since replacing the whole set must never touch an approved account's kids.
func (s *Service) ReplaceKids(ctx context.Context, req ReplaceKidsRequest) ([]*Kid, error) {
	userID := ctx.Value(s.cfg.JWT.User.Key).(*Claims).ID

	names := make([]string, 0, len(req.FirstNames))
	for _, n := range req.FirstNames {
		if trimmed := strings.TrimSpace(n); trimmed != "" {
			names = append(names, trimmed)
		}
	}

	kids, err := s.store.replaceKids(ctx, userID, names)
	if err != nil {
		return nil, err
	}
	if kids == nil {
		kids = []*Kid{}
	}
	return kids, nil
}

func (s *Service) UpdateKid(ctx context.Context, id string, req KidRequest) (*Kid, error) {
	userID := ctx.Value(s.cfg.JWT.User.Key).(*Claims).ID
	kidID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperrors.NewBadRequest("invalid kid id", err)
	}
	return s.store.updateKid(ctx, kidID, userID, req.FirstName, req.LastName)
}

func (s *Service) RemoveKid(ctx context.Context, id string) error {
	userID := ctx.Value(s.cfg.JWT.User.Key).(*Claims).ID
	kidID, err := uuid.Parse(id)
	if err != nil {
		return apperrors.NewBadRequest("invalid kid id", err)
	}
	return s.store.deleteKid(ctx, kidID, userID)
}

func (s *Service) UpdateKidPicture(ctx context.Context, id string, contentType string, data io.Reader) (*Kid, error) {
	userID := ctx.Value(s.cfg.JWT.User.Key).(*Claims).ID
	kidID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperrors.NewBadRequest("invalid kid id", err)
	}

	// Verify ownership before uploading so we never store an orphan object for a
	// kid the caller does not own.
	if _, err := s.store.getOwnedKid(ctx, kidID, userID); err != nil {
		return nil, err
	}

	url, err := s.storageClient.UploadPicture(ctx, contentType, data)
	if err != nil {
		return nil, err
	}

	return s.store.updateKidPicture(ctx, kidID, userID, url)
}
