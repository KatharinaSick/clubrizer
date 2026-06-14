package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/katharinasick/clubrizer/internal/api"
	"github.com/katharinasick/clubrizer/internal/app"
	"github.com/katharinasick/clubrizer/internal/email"
	"github.com/katharinasick/clubrizer/internal/events"
	"github.com/katharinasick/clubrizer/internal/rbac"
	"github.com/katharinasick/clubrizer/internal/storage"
	"github.com/katharinasick/clubrizer/internal/users"
)

func main() {
	log := app.NewLogger()

	cfg, err := app.ParseConfig()
	if err != nil {
		log.Error("Failed to parse config. Terminating.", "error", err)
		os.Exit(1)
	}

	dbPool, err := app.ConnectDatastore(log, cfg)
	if err != nil {
		log.Error("Failed to connect to database. Terminating.", "error", err)
		os.Exit(1)
	}

	emailClient := email.NewClient(cfg.Resend.ApiKey, cfg.Resend.FromAddress)

	storageClient, err := storage.NewClient(context.Background(), cfg.GCS.ProfilePicturesBucketName)
	if err != nil {
		log.Error("Failed to create GCS client. Terminating.", "error", err)
		os.Exit(1)
	}

	httpServer := &http.Server{
		Addr: net.JoinHostPort("", cfg.Server.Port),
		Handler: api.NewHandler(
			*cfg,
			users.NewService(log, cfg, dbPool, emailClient, storageClient),
			events.NewService(log, cfg, dbPool, rbac.NewService(dbPool)),
		),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Info("Starting server", "address", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Error("Failed to start server. Terminating.", "error", err)
		os.Exit(1)
	}
}
