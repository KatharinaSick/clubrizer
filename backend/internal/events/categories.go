package events

import (
	"context"

	"github.com/katharinasick/clubrizer/internal/users"
)

func (s *Service) ListCategories(ctx context.Context) ([]*Category, error) {
	categories, err := s.store.getAllCategories(ctx)
	if err != nil {
		return nil, err
	}

	userID := ctx.Value(s.cfg.JWT.User.Key).(*users.Claims).ID
	creatableIDs, err := s.rbac.GetCreatableCategoryIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, c := range categories {
		canCreate := creatableIDs[c.ID]
		c.CanCreate = &canCreate
	}

	return categories, nil
}
