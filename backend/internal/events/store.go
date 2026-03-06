package events

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
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
	rows, err := s.conn.Query(ctx, "SELECT * FROM event_categories ORDER BY sort_order ASC")
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

func (s *store) getFutureEvents(ctx context.Context) ([]*Event, error) {
	rows, err := s.conn.Query(ctx, `
		SELECT 
			e.id, e.title, e.description, e.location, e.start_time, e.created_by, e.created_at, e.category,
			c.id, c.name, c.color, c.picture, c.sort_order
		FROM events e
		LEFT JOIN event_categories c ON e.category = c.id
		WHERE e.start_time >= NOW()
		ORDER BY e.start_time
	`)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query events: %s", err.Error()))
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		var e Event
		var c Category
		// We don't need created_by/at for category in the list view usually, or we can scan them if needed.
		// For now, scanning the main fields.
		err := rows.Scan(
			&e.ID, &e.Title, &e.Description, &e.Location, &e.StartTime, &e.CreatedBy, &e.CreatedAt, &e.CategoryID,
			&c.ID, &c.Name, &c.Color, &c.Picture, &c.SortOrder,
		)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to scan event: %s", err.Error()))
		}
		e.Category = c
		events = append(events, &e)
	}

	if rows.Err() != nil {
		return nil, errors.New(fmt.Sprintf("rows error: %s", rows.Err().Error()))
	}

	return events, nil
}

func (s *store) createEvent(ctx context.Context, e *Event) (uuid.UUID, error) {
	var id uuid.UUID
	err := s.conn.QueryRow(
		context.Background(),
		"INSERT INTO events(title, category, description, location, start_time, created_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		e.Title, e.CategoryID, e.Description, e.Location, e.StartTime, ctx.Value(s.cfg.OAuth.User.Key).(*users.Claims).ID,
	).Scan(&id)

	if err != nil {
		return uuid.Nil, errors.New(fmt.Sprintf("failed to create event: %s", err.Error()))
	}

	return id, nil
}
