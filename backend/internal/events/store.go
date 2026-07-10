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
			e.id, e.title, e.description, e.location, e.start_time, e.created_by, e.created_at, e.category, e.cancelled_at,
			c.id, c.name, c.color, c.picture, c.sort_order, c.custom_label,
			u.id, u.given_name, u.family_name, COALESCE(u.nick_name, u.given_name), u.picture
		FROM events e
		LEFT JOIN event_categories c ON e.category = c.id
		LEFT JOIN users u ON e.created_by = u.id
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
		var cr Creator
		err := rows.Scan(
			&e.ID, &e.Title, &e.Description, &e.Location, &e.StartTime, &e.CreatedBy, &e.CreatedAt, &e.CategoryID, &e.CancelledAt,
			&c.ID, &c.Name, &c.Color, &c.Picture, &c.SortOrder, &c.CustomLabel,
			&cr.ID, &cr.GivenName, &cr.FamilyName, &cr.NickName, &cr.Picture,
		)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to scan event: %s", err.Error()))
		}
		e.Category = c
		e.Creator = cr
		events = append(events, &e)
	}

	if rows.Err() != nil {
		return nil, errors.New(fmt.Sprintf("rows error: %s", rows.Err().Error()))
	}

	return events, nil
}

func (s *store) getEventById(ctx context.Context, id uuid.UUID) (*Event, error) {
	row := s.conn.QueryRow(ctx, `
		SELECT
			e.id, e.title, e.description, e.location, e.start_time, e.created_by, e.created_at, e.category, e.cancelled_at,
			c.id, c.name, c.color, c.picture, c.sort_order, c.custom_label,
			u.id, u.given_name, u.family_name, COALESCE(u.nick_name, u.given_name), u.picture
		FROM events e
		LEFT JOIN event_categories c ON e.category = c.id
		LEFT JOIN users u ON e.created_by = u.id
		WHERE e.id = $1
	`, id)

	var e Event
	var c Category
	var cr Creator
	err := row.Scan(
		&e.ID, &e.Title, &e.Description, &e.Location, &e.StartTime, &e.CreatedBy, &e.CreatedAt, &e.CategoryID, &e.CancelledAt,
		&c.ID, &c.Name, &c.Color, &c.Picture, &c.SortOrder, &c.CustomLabel,
		&cr.ID, &cr.GivenName, &cr.FamilyName, &cr.NickName, &cr.Picture,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewNotFound(fmt.Sprintf("event with id %s not found", id))
		}
		return nil, errors.New(fmt.Sprintf("failed to scan event: %s", err.Error()))
	}
	e.Category = c
	e.Creator = cr

	return &e, nil
}

func (s *store) getEventResponses(ctx context.Context, eventId uuid.UUID, userId uuid.UUID) (*EventResponses, error) {
	rows, err := s.conn.Query(ctx, `
		SELECT u.id, u.given_name, u.family_name, COALESCE(u.nick_name, u.given_name), u.picture, r.response
		FROM event_responses r
		JOIN users u ON r.user_id = u.id
		WHERE r.event_id = $1
		ORDER BY r.created_at
	`, eventId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query event responses: %s", err.Error()))
	}
	defer rows.Close()

	responses := &EventResponses{Attendees: []*EventAttendee{}}
	for rows.Next() {
		var a EventAttendee
		err := rows.Scan(&a.ID, &a.GivenName, &a.FamilyName, &a.NickName, &a.Picture, &a.Response)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to scan event response: %s", err.Error()))
		}
		if a.Response {
			responses.Going++
		} else {
			responses.NotGoing++
		}
		if a.ID == userId {
			r := a.Response
			responses.CurrentUserResponse = &r
		}
		responses.Attendees = append(responses.Attendees, &a)
	}

	if rows.Err() != nil {
		return nil, errors.New(fmt.Sprintf("rows error: %s", rows.Err().Error()))
	}

	return responses, nil
}

func (s *store) upsertEventResponse(ctx context.Context, eventId uuid.UUID, userId uuid.UUID, response bool) error {
	_, err := s.conn.Exec(ctx, `
		INSERT INTO event_responses (event_id, user_id, response)
		VALUES ($1, $2, $3)
		ON CONFLICT (event_id, user_id) DO UPDATE SET response = EXCLUDED.response
	`, eventId, userId, response)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to upsert event response: %s", err.Error()))
	}
	return nil
}

func (s *store) deleteEvent(ctx context.Context, id uuid.UUID) error {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to begin transaction: %s", err.Error()))
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, "DELETE FROM event_responses WHERE event_id = $1", id); err != nil {
		return errors.New(fmt.Sprintf("failed to delete event responses: %s", err.Error()))
	}
	if _, err := tx.Exec(ctx, "DELETE FROM events WHERE id = $1", id); err != nil {
		return errors.New(fmt.Sprintf("failed to delete event: %s", err.Error()))
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.New(fmt.Sprintf("failed to commit transaction: %s", err.Error()))
	}
	return nil
}

func (s *store) cancelEvent(ctx context.Context, id uuid.UUID) error {
	_, err := s.conn.Exec(ctx, "UPDATE events SET cancelled_at = NOW() WHERE id = $1", id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to cancel event: %s", err.Error()))
	}
	return nil
}

func (s *store) uncancelEvent(ctx context.Context, id uuid.UUID) error {
	_, err := s.conn.Exec(ctx, "UPDATE events SET cancelled_at = NULL WHERE id = $1", id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to uncancel event: %s", err.Error()))
	}
	return nil
}

func (s *store) createEvent(ctx context.Context, e *Event) (uuid.UUID, error) {
	var id uuid.UUID
	err := s.conn.QueryRow(
		context.Background(),
		"INSERT INTO events(title, category, description, location, start_time, created_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		e.Title, e.CategoryID, e.Description, e.Location, e.StartTime, ctx.Value(s.cfg.JWT.User.Key).(*users.Claims).ID,
	).Scan(&id)

	if err != nil {
		return uuid.Nil, errors.New(fmt.Sprintf("failed to create event: %s", err.Error()))
	}

	return id, nil
}
