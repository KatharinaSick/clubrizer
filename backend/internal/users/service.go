package users

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katharinasick/clubrizer/internal/app"
	"github.com/katharinasick/clubrizer/internal/email"
	"github.com/katharinasick/clubrizer/internal/storage"
)

// Service is responsible for everything related to handling users on the server. This includes operations like authentication, approvals, crud,...
type Service struct {
	log           app.Logger
	cfg           *app.Config
	store         *store
	email         *email.Client
	storageClient *storage.Client
}

func NewService(log app.Logger, cfg *app.Config, conn *pgxpool.Pool, emailClient *email.Client, storageClient *storage.Client) *Service {
	store := newStore(log, cfg, conn)
	return &Service{log, cfg, store, emailClient, storageClient}
}
