package main

import (
	"github.com/katharinasick/clubrizer/internal/api"
	"github.com/katharinasick/clubrizer/internal/setup"
	"github.com/katharinasick/clubrizer/internal/users"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	log := setup.NewLogger()

	cfg, err := setup.ParseConfig()
	if err != nil {
		log.Error("Failed to parse config. Terminating.", "error", err)
		os.Exit(1)
	}

	dbPool, err := setup.ConnectDatastore(log, cfg)
	if err != nil {
		log.Error("Failed to connect to database. Terminating.", "error", err)
		os.Exit(1)
	}

	httpServer := &http.Server{
		Addr: net.JoinHostPort("", cfg.Server.Port),
		Handler: api.NewHandler(
			users.NewService(log, cfg, dbPool),
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
