package rbac

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	store *store
}

func NewService(conn *pgxpool.Pool) *Service {
	return &Service{store: newStore(conn)}
}

// IsAuthorized returns true if the user has the given permission, or is an admin (who bypasses all checks).
func (s *Service) IsAuthorized(ctx context.Context, userID uuid.UUID, permissionKey string) (bool, error) {
	isAdmin, err := s.store.hasAdminRole(ctx, userID)
	if err != nil {
		return false, err
	}
	if isAdmin {
		return true, nil
	}
	return s.store.hasPermission(ctx, userID, permissionKey)
}

// IsAuthorizedToCreateEvent returns true if the user holds a role that is allowed to create events
// in the given category, or is an admin.
func (s *Service) IsAuthorizedToCreateEvent(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) (bool, error) {
	isAdmin, err := s.store.hasAdminRole(ctx, userID)
	if err != nil {
		return false, err
	}
	if isAdmin {
		return true, nil
	}
	return s.store.canCreateInCategory(ctx, userID, categoryID)
}

// GetCreatableCategoryIDs returns a map of all category IDs the user may create events in.
func (s *Service) GetCreatableCategoryIDs(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]bool, error) {
	isAdmin, err := s.store.hasAdminRole(ctx, userID)
	if err != nil {
		return nil, err
	}

	var ids []uuid.UUID
	if isAdmin {
		ids, err = s.store.getAllCategoryIDs(ctx)
	} else {
		ids, err = s.store.getCreatableCategoryIDs(ctx, userID)
	}

	ids, err := s.store.getCreatableCategoryIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make(map[uuid.UUID]bool, len(ids))
	for _, id := range ids {
		result[id] = true
	}
	return result, nil
}
