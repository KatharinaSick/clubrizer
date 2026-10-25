package notifications

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
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

func (s *store) touchSubscription(ctx context.Context, id uuid.UUID) error {
	_, err := s.conn.Exec(ctx, "UPDATE push_subscriptions SET last_used_at = current_timestamp WHERE id = $1", id)
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

func (s *store) getSubscriptionsByUser(ctx context.Context, userId uuid.UUID) ([]*Subscription, error) {
	rows, err := s.conn.Query(ctx, "SELECT id, user_id, endpoint, p256dh, auth FROM push_subscriptions WHERE user_id = $1", userId)
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
