package main

import (
	"github.com/katharinasick/clubrizer/internal/api"
	"github.com/katharinasick/clubrizer/internal/app"
	"github.com/katharinasick/clubrizer/internal/events"
	"github.com/katharinasick/clubrizer/internal/users"
	"net"
	"net/http"
	"os"
	"time"
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

	httpServer := &http.Server{
		Addr: net.JoinHostPort("", cfg.Server.Port),
		Handler: api.NewHandler(
			*cfg,
			users.NewService(log, cfg, dbPool),
			events.NewService(log, cfg, dbPool),
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
