package events

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katharinasick/clubrizer/internal/app"
)

// Service is responsible for everything related to handling events.
type Service struct {
	log   app.Logger
	cfg   *app.Config
	store *store
}

func NewService(log app.Logger, cfg *app.Config, conn *pgxpool.Pool) *Service {
	store := newStore(log, cfg, conn)
	return &Service{log, cfg, store}
}
