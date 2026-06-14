package events

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katharinasick/clubrizer/internal/app"
)

type rbacService interface {
	IsAuthorizedToCreateEvent(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) (bool, error)
	GetCreatableCategoryIDs(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]bool, error)
}

// Service is responsible for everything related to handling events.
type Service struct {
	log   app.Logger
	cfg   *app.Config
	store *store
	rbac  rbacService
}

func NewService(log app.Logger, cfg *app.Config, conn *pgxpool.Pool, rbac rbacService) *Service {
	store := newStore(log, cfg, conn)
	return &Service{log, cfg, store, rbac}
}
