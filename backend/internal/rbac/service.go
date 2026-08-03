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

// GetEffectivePermissions returns the permission keys the user effectively holds. Admins
// bypass individual permission checks, so they expand to the full catalog (AllPermissions);
// everyone else gets exactly the keys assigned to their roles. Intended for surfacing
// permissions to the frontend (e.g. in JWT claims) so it can show or hide UI. Backend
// authorization must still go through IsAuthorized, which reads the DB as the source of
// truth — a stale token must never grant access.
func (s *Service) GetEffectivePermissions(ctx context.Context, userID uuid.UUID) ([]string, error) {
	isAdmin, err := s.store.hasAdminRole(ctx, userID)
	if err != nil {
		return nil, err
	}
	if isAdmin {
		// Return a copy so callers can't mutate the shared catalog slice.
		return append([]string(nil), AllPermissions...), nil
	}

	keys, err := s.store.getPermissionKeys(ctx, userID)
	if err != nil {
		return nil, err
	}
	if keys == nil {
		keys = []string{}
	}
	return keys, nil
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

	if err != nil {
		return nil, err
	}

	result := make(map[uuid.UUID]bool, len(ids))
	for _, id := range ids {
		result[id] = true
	}
	return result, nil
}
