package app

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

func ConnectDatastore(log Logger, cfg *Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.Database.Url)
	if err != nil {
		log.Error("Failed to parse database config. Terminating.", "error", err)
		os.Exit(1)
	}
	poolCfg.MaxConns = 10
	poolCfg.MinConns = 2

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		log.Error("Failed to connect to database. Terminating.", "error", err)
		os.Exit(1)
	}

	// pool.Close() is intentionally omitted — Cloud Run sends SIGTERM and terminates the process, which closes connections automatically.
	return pool, nil
}
