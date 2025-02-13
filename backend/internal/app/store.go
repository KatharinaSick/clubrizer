package app

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

func ConnectDatastore(log Logger, cfg *Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), cfg.Database.Url)
	if err != nil {
		log.Error("Failed to connect to database. Terminating.", "error", err)
		os.Exit(1)
	}

	// TODO
	// defer func(pool *pgxpool.Pool, ctx context.Context) {
	//pool.Close()
	//}(pool, context.Background())
	return pool, nil
}
