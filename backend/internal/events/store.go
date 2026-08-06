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
		-- Keep an event on the list until 4 hours after it starts, so an event in progress or
		-- just finished stays visible. A plain interval avoids any timezone/day-boundary logic.
		WHERE e.start_time >= NOW() - INTERVAL '4 hours'
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

func (s *store) getEventResponses(ctx context.Context, eventId uuid.UUID, r responder) (*EventResponses, error) {
	// Attendee grid: every response for the event, whether it belongs to an account
	// holder (user_id) or a kid (kid_id). COALESCE folds kid fields over user fields
	// so both render the same way; for a kid, the parent's nick name is joined in.
	rows, err := s.conn.Query(ctx, `
		SELECT
			COALESCE(u.id, k.id),
			COALESCE(u.given_name, k.given_name, ''),
			COALESCE(u.family_name, k.family_name, ''),
			COALESCE(u.nick_name, u.given_name, k.given_name, ''),
			COALESCE(u.picture, k.picture),
			resp.response,
			(resp.kid_id IS NOT NULL),
			pu.nick_name
		FROM event_responses resp
		LEFT JOIN users u ON resp.user_id = u.id
		LEFT JOIN kids k ON resp.kid_id = k.id
		LEFT JOIN users pu ON k.user_id = pu.id
		WHERE resp.event_id = $1
		  -- Own responses always show; a kid response only if the kid isn't rejected.
		  -- Removed (soft-deleted) kids keep status 'approved', so their past responses
		  -- stay visible here even though they're gone from the parent's controls.
		  AND (resp.kid_id IS NULL OR k.status <> 'rejected')
		ORDER BY resp.created_at
	`, eventId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query event responses: %s", err.Error()))
	}
	defer rows.Close()

	responses := &EventResponses{Attendees: []*EventAttendee{}, MyResponses: []*MyResponse{}}
	var ownResponse *bool
	for rows.Next() {
		var a EventAttendee
		var parent *string
		err := rows.Scan(&a.ID, &a.GivenName, &a.FamilyName, &a.NickName, &a.Picture, &a.Response, &a.IsKid, &parent)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to scan event response: %s", err.Error()))
		}
		if a.IsKid {
			a.Parent = parent
		}
		if a.Response {
			responses.Going++
		} else {
			responses.NotGoing++
		}
		if !a.IsKid && a.ID == r.UserID {
			resp := a.Response
			ownResponse = &resp
		}
		responses.Attendees = append(responses.Attendees, &a)
	}
	if rows.Err() != nil {
		return nil, errors.New(fmt.Sprintf("rows error: %s", rows.Err().Error()))
	}

	// The account holder participates only if self_participates; their own response
	// was captured above (nil if they haven't responded yet).
	if r.SelfParticipates {
		responses.MyResponses = append(responses.MyResponses, &MyResponse{
			ID:       r.UserID,
			IsSelf:   true,
			Name:     r.Name,
			Picture:  r.Picture,
			Response: ownResponse,
		})
	}

	// Approved kids the account manages, each with their response for this event.
	kidRows, err := s.conn.Query(ctx, `
		SELECT k.id, COALESCE(k.given_name, ''), k.picture, resp.response
		FROM kids k
		LEFT JOIN event_responses resp ON resp.kid_id = k.id AND resp.event_id = $2
		WHERE k.user_id = $1 AND k.status = 'approved' AND k.deleted_at IS NULL
		ORDER BY k.created_at
	`, r.UserID, eventId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query kid responses: %s", err.Error()))
	}
	defer kidRows.Close()

	for kidRows.Next() {
		var m MyResponse
		if err := kidRows.Scan(&m.ID, &m.Name, &m.Picture, &m.Response); err != nil {
			return nil, errors.New(fmt.Sprintf("failed to scan kid response: %s", err.Error()))
		}
		responses.MyResponses = append(responses.MyResponses, &m)
	}
	if kidRows.Err() != nil {
		return nil, errors.New(fmt.Sprintf("rows error: %s", kidRows.Err().Error()))
	}

	return responses, nil
}

func (s *store) upsertEventResponse(ctx context.Context, eventId uuid.UUID, userId uuid.UUID, response bool) error {
	// The WHERE predicate is required to match the unique index on (event_id, user_id):
	// since user_id became nullable (migration 000010, for kid responses), that index is
	// partial — WHERE user_id IS NOT NULL — on CockroachDB, which rewrites a unique index
	// over a nullable column as partial. Without the predicate the arbiter doesn't match
	// (SQLSTATE 42P10). Mirrors the kid upsert. Safe on Postgres too: a predicate may infer
	// a non-partial index there.
	_, err := s.conn.Exec(ctx, `
		INSERT INTO event_responses (event_id, user_id, response)
		VALUES ($1, $2, $3)
		ON CONFLICT (event_id, user_id) WHERE user_id IS NOT NULL DO UPDATE SET response = EXCLUDED.response
	`, eventId, userId, response)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to upsert event response: %s", err.Error()))
	}
	return nil
}

func (s *store) upsertKidEventResponse(ctx context.Context, eventId uuid.UUID, kidId uuid.UUID, response bool) error {
	_, err := s.conn.Exec(ctx, `
		INSERT INTO event_responses (event_id, kid_id, response)
		VALUES ($1, $2, $3)
		ON CONFLICT (event_id, kid_id) WHERE kid_id IS NOT NULL DO UPDATE SET response = EXCLUDED.response
	`, eventId, kidId, response)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to upsert kid event response: %s", err.Error()))
	}
	return nil
}

// getOwnedKidStatus returns the kid's approval status only if it belongs to userId and
// hasn't been removed. A kid owned by someone else — or one that was soft-deleted — is
// reported as not found, so a guessed id reveals nothing and a removed kid can't receive
// new responses.
func (s *store) getOwnedKidStatus(ctx context.Context, kidId uuid.UUID, userId uuid.UUID) (string, error) {
	var status string
	err := s.conn.QueryRow(ctx,
		"SELECT status FROM kids WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL",
		kidId, userId,
	).Scan(&status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", apperrors.NewNotFound(fmt.Sprintf("kid with id %s not found", kidId))
		}
		return "", errors.New(fmt.Sprintf("failed to query kid status: %s", err.Error()))
	}
	return status, nil
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
		ctx,
		"INSERT INTO events(title, category, description, location, start_time, created_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		e.Title, e.CategoryID, e.Description, e.Location, e.StartTime, ctx.Value(s.cfg.JWT.User.Key).(*users.Claims).ID,
	).Scan(&id)

	if err != nil {
		return uuid.Nil, errors.New(fmt.Sprintf("failed to create event: %s", err.Error()))
	}

	return id, nil
}

func (s *store) getEventComments(ctx context.Context, eventId uuid.UUID) ([]*Comment, error) {
	rows, err := s.conn.Query(ctx, `
		SELECT
			cm.id, cm.body, cm.created_at,
			u.id, u.given_name, u.family_name, COALESCE(u.nick_name, u.given_name), u.picture
		FROM event_comments cm
		LEFT JOIN users u ON cm.user_id = u.id
		WHERE cm.event_id = $1
		ORDER BY cm.created_at ASC
	`, eventId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query event comments: %s", err.Error()))
	}
	defer rows.Close()

	comments := []*Comment{}
	for rows.Next() {
		var c Comment
		var a Creator
		if err := rows.Scan(&c.ID, &c.Body, &c.CreatedAt, &a.ID, &a.GivenName, &a.FamilyName, &a.NickName, &a.Picture); err != nil {
			return nil, errors.New(fmt.Sprintf("failed to scan event comment: %s", err.Error()))
		}
		c.Author = a
		comments = append(comments, &c)
	}
	if rows.Err() != nil {
		return nil, errors.New(fmt.Sprintf("rows error: %s", rows.Err().Error()))
	}

	return comments, nil
}

func (s *store) createComment(ctx context.Context, eventId uuid.UUID, userId uuid.UUID, body string) (*Comment, error) {
	row := s.conn.QueryRow(ctx, `
		WITH inserted AS (
			INSERT INTO event_comments (event_id, user_id, body)
			VALUES ($1, $2, $3)
			RETURNING id, body, created_at, user_id
		)
		SELECT
			i.id, i.body, i.created_at,
			u.id, u.given_name, u.family_name, COALESCE(u.nick_name, u.given_name), u.picture
		FROM inserted i
		LEFT JOIN users u ON i.user_id = u.id
	`, eventId, userId, body)

	var c Comment
	var a Creator
	if err := row.Scan(&c.ID, &c.Body, &c.CreatedAt, &a.ID, &a.GivenName, &a.FamilyName, &a.NickName, &a.Picture); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to create comment: %s", err.Error()))
	}
	c.Author = a

	return &c, nil
}
