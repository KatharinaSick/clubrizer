package events

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katharinasick/clubrizer/internal/app"
	"github.com/katharinasick/clubrizer/internal/apperrors"
	"github.com/katharinasick/clubrizer/internal/users"
)

type store struct {
	log  app.Logger
	cfg  *app.Config
	conn *pgxpool.Pool
}

func newStore(log app.Logger, cfg *app.Config, conn *pgxpool.Pool) *store {
	return &store{log, cfg, conn}
}

func (s *store) getAllCategories(ctx context.Context) ([]*Category, error) {
	rows, err := s.conn.Query(ctx, "SELECT * FROM event_categories")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query event categories: %s", err.Error()))
	}

	c, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Category])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewNotFound(fmt.Sprintf("no categories found"))
		}
		return nil, errors.New(fmt.Sprintf("failed to scan categories: %s", err.Error()))
	}

	return c, nil
}

func (s *store) getAllEvents(ctx context.Context) ([]*Event, error) {
	rows, err := s.conn.Query(ctx, "SELECT * from events e "+
		"LEFT JOIN event_categories c ON e.category = c.id "+
		"ORDER BY e.start_time",
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query events: %s", err.Error()))
	}

	c, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Event])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewNotFound(fmt.Sprintf("no events found"))
		}
		return nil, errors.New(fmt.Sprintf("failed to scan events: %s", err.Error()))
	}

	return c, nil
}

func (s *store) createEvent(ctx context.Context, e *Event) error {
	_, err := s.conn.Exec(
		context.Background(),
		"INSERT INTO events(title, category, description, location, start_time, created_by) VALUES ($1, $2, $3, $4, $5, $6)",
		e.Title, e.Category, e.Description, e.Location, e.StartTime, ctx.Value(s.cfg.OAuth.User.Key).(*users.Claims).ID,
	)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create event: %s", err.Error()))
	}

	return nil
}
