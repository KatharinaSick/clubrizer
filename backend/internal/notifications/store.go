package notifications

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katharinasick/clubrizer/internal/app"
)

type store struct {
	log  app.Logger
	cfg  *app.Config
	conn *pgxpool.Pool
}

func newStore(log app.Logger, cfg *app.Config, conn *pgxpool.Pool) *store {
	return &store{log, cfg, conn}
}

// upsertSubscription stores the subscription for the user. An endpoint belongs to exactly one
// browser, so if it is already known (e.g. someone else was signed in on this device before)
// it is reassigned to the current user — the device should only get the current user's
// notifications.
func (s *store) upsertSubscription(ctx context.Context, userId uuid.UUID, endpoint, p256dh, auth string) error {
	_, err := s.conn.Exec(ctx, `
		INSERT INTO push_subscriptions (user_id, endpoint, p256dh, auth)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (endpoint) DO UPDATE
		SET user_id = excluded.user_id, p256dh = excluded.p256dh, auth = excluded.auth, last_used_at = current_timestamp
	`, userId, endpoint, p256dh, auth)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to upsert push subscription: %s", err.Error()))
	}
	return nil
}

// deleteSubscription removes one of the user's subscriptions. Scoped to the user so nobody
// can remove another user's device.
func (s *store) deleteSubscription(ctx context.Context, userId uuid.UUID, endpoint string) error {
	_, err := s.conn.Exec(ctx, "DELETE FROM push_subscriptions WHERE user_id = $1 AND endpoint = $2", userId, endpoint)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete push subscription: %s", err.Error()))
	}
	return nil
}

// deleteSubscriptionOfDevice removes the subscription with this endpoint and secret, no matter
// which user it belongs to — matching the secret proves the caller is on that device.
func (s *store) deleteSubscriptionOfDevice(ctx context.Context, endpoint, auth string) error {
	_, err := s.conn.Exec(ctx, "DELETE FROM push_subscriptions WHERE endpoint = $1 AND auth = $2", endpoint, auth)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete push subscription: %s", err.Error()))
	}
	return nil
}

func (s *store) deleteAllSubscriptions(ctx context.Context, userId uuid.UUID) error {
	_, err := s.conn.Exec(ctx, "DELETE FROM push_subscriptions WHERE user_id = $1", userId)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete push subscriptions: %s", err.Error()))
	}
	return nil
}

// deleteSubscriptionById removes a subscription the push service reported as gone.
func (s *store) deleteSubscriptionById(ctx context.Context, id uuid.UUID) error {
	_, err := s.conn.Exec(ctx, "DELETE FROM push_subscriptions WHERE id = $1", id)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete push subscription: %s", err.Error()))
	}
	return nil
}

// touchSubscriptions records that the push service accepted a message for these devices.
func (s *store) touchSubscriptions(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	_, err := s.conn.Exec(ctx, "UPDATE push_subscriptions SET last_used_at = current_timestamp WHERE id = ANY($1)", ids)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to update push subscription: %s", err.Error()))
	}
	return nil
}

func (s *store) hasSubscription(ctx context.Context, userId uuid.UUID, endpoint string) (bool, error) {
	var exists bool
	err := s.conn.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM push_subscriptions WHERE user_id = $1 AND endpoint = $2)", userId, endpoint).Scan(&exists)
	if err != nil {
		return false, errors.New(fmt.Sprintf("failed to check push subscription: %s", err.Error()))
	}
	return exists, nil
}

// getSubscriptionsByUsers returns every device of the given users. Duplicate ids are fine —
// each device is returned once.
func (s *store) getSubscriptionsByUsers(ctx context.Context, userIds []uuid.UUID) ([]*Subscription, error) {
	rows, err := s.conn.Query(ctx, "SELECT id, user_id, endpoint, p256dh, auth FROM push_subscriptions WHERE user_id = ANY($1)", userIds)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query push subscriptions: %s", err.Error()))
	}
	defer rows.Close()

	var subs []*Subscription
	for rows.Next() {
		var sub Subscription
		if err := rows.Scan(&sub.ID, &sub.UserID, &sub.Endpoint, &sub.P256dh, &sub.Auth); err != nil {
			return nil, errors.New(fmt.Sprintf("failed to scan push subscription: %s", err.Error()))
		}
		subs = append(subs, &sub)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to read push subscriptions: %s", err.Error()))
	}

	return subs, nil
}

// getEmails returns the email addresses of the given users.
func (s *store) getEmails(ctx context.Context, userIds []uuid.UUID) ([]string, error) {
	rows, err := s.conn.Query(ctx, "SELECT email FROM users WHERE id = ANY($1)", userIds)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query emails: %s", err.Error()))
	}
	emails, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to collect emails: %s", err.Error()))
	}
	return emails, nil
}

// getApprovedUserIds returns the ids of all approved accounts.
func (s *store) getApprovedUserIds(ctx context.Context) ([]uuid.UUID, error) {
	rows, err := s.conn.Query(ctx, "SELECT id FROM users WHERE status = 'approved'")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query approved users: %s", err.Error()))
	}
	ids, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to collect approved users: %s", err.Error()))
	}
	return ids, nil
}

// getOptOuts returns the categories the user switched off.
func (s *store) getOptOuts(ctx context.Context, userId uuid.UUID) (map[string]bool, error) {
	rows, err := s.conn.Query(ctx, "SELECT category FROM notification_opt_outs WHERE user_id = $1", userId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query notification opt-outs: %s", err.Error()))
	}
	categories, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to collect notification opt-outs: %s", err.Error()))
	}
	optOuts := make(map[string]bool, len(categories))
	for _, c := range categories {
		optOuts[c] = true
	}
	return optOuts, nil
}

// getOptedOutUserIds returns which of the given users switched the category off.
func (s *store) getOptedOutUserIds(ctx context.Context, category string, userIds []uuid.UUID) (map[uuid.UUID]bool, error) {
	rows, err := s.conn.Query(ctx, "SELECT user_id FROM notification_opt_outs WHERE category = $1 AND user_id = ANY($2)", category, userIds)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query notification opt-outs: %s", err.Error()))
	}
	ids, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to collect notification opt-outs: %s", err.Error()))
	}
	optedOut := make(map[uuid.UUID]bool, len(ids))
	for _, id := range ids {
		optedOut[id] = true
	}
	return optedOut, nil
}

func (s *store) setOptOut(ctx context.Context, userId uuid.UUID, category string, optOut bool) error {
	query := "DELETE FROM notification_opt_outs WHERE user_id = $1 AND category = $2"
	if optOut {
		query = "INSERT INTO notification_opt_outs (user_id, category) VALUES ($1, $2) ON CONFLICT (user_id, category) DO NOTHING"
	}
	if _, err := s.conn.Exec(ctx, query, userId, category); err != nil {
		return errors.New(fmt.Sprintf("failed to update notification opt-out: %s", err.Error()))
	}
	return nil
}
