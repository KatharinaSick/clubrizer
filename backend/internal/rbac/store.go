package rbac

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type store struct {
	conn *pgxpool.Pool
}

func newStore(conn *pgxpool.Pool) *store {
	return &store{conn: conn}
}

func (s *store) hasAdminRole(ctx context.Context, userID uuid.UUID) (bool, error) {
	var exists bool
	err := s.conn.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM user_roles ur
			JOIN roles r ON r.id = ur.role_id
			WHERE ur.user_id = $1 AND r.name = 'admin'
		)
	`, userID).Scan(&exists)
	return exists, err
}

func (s *store) hasPermission(ctx context.Context, userID uuid.UUID, permissionKey string) (bool, error) {
	var exists bool
	err := s.conn.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM user_roles ur
			JOIN role_permissions rp ON rp.role_id = ur.role_id
			WHERE ur.user_id = $1 AND rp.permission_key = $2
		)
	`, userID, permissionKey).Scan(&exists)
	return exists, err
}

func (s *store) canCreateInCategory(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) (bool, error) {
	var exists bool
	err := s.conn.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM user_roles ur
			JOIN category_creator_roles ccr ON ccr.role_id = ur.role_id
			WHERE ur.user_id = $1 AND ccr.category_id = $2
		)
	`, userID, categoryID).Scan(&exists)
	return exists, err
}

func (s *store) getCreatableCategoryIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := s.conn.Query(ctx, `
		SELECT DISTINCT ccr.category_id
		FROM user_roles ur
		JOIN category_creator_roles ccr ON ccr.role_id = ur.role_id
		WHERE ur.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (s *store) getAllCategoryIDs(ctx context.Context) ([]uuid.UUID, error) {
	rows, err := s.conn.Query(ctx, `SELECT id FROM event_categories`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
