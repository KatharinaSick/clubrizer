package users

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katharinasick/clubrizer/internal/app"
	"github.com/katharinasick/clubrizer/internal/apperrors"
)

type store struct {
	log  app.Logger
	cfg  *app.Config
	conn *pgxpool.Pool
}

func newStore(log app.Logger, cfg *app.Config, conn *pgxpool.Pool) *store {
	return &store{log, cfg, conn}
}

func (s *store) getUserById(ctx context.Context, id string) (*User, error) {
	rows, err := s.conn.Query(ctx, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query user: %s", err.Error()))
	}

	u, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewNotFound(fmt.Sprintf("user with id %s not found", id))
		}
		return nil, errors.New(fmt.Sprintf("failed to scan user: %s", err.Error()))
	}

	return u, nil
}

func (s *store) getUserByMail(ctx context.Context, email string) (*User, error) {
	rows, err := s.conn.Query(ctx, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query user: %s", err.Error()))
	}

	u, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewNotFound(fmt.Sprintf("user with email %s not found", email))
		}
		return nil, errors.New(fmt.Sprintf("failed to scan user: %s", err.Error()))
	}

	return u, nil
}

func (s *store) countRecentOTPRequests(ctx context.Context, email string) (int, error) {
	var count int
	err := s.conn.QueryRow(ctx,
		"SELECT COUNT(*) FROM otp_tokens WHERE email = $1 AND created_at > NOW() - INTERVAL '1 hour'",
		email,
	).Scan(&count)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("failed to count OTP requests: %s", err.Error()))
	}
	return count, nil
}

func (s *store) invalidateActiveOTPs(ctx context.Context, email string) error {
	_, err := s.conn.Exec(ctx,
		"UPDATE otp_tokens SET invalidated_at = NOW() WHERE email = $1 AND invalidated_at IS NULL AND expires_at > NOW()",
		email,
	)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to invalidate OTP tokens: %s", err.Error()))
	}
	return nil
}

func (s *store) getActiveOTPByEmail(ctx context.Context, email string) (*OTPToken, error) {
	rows, err := s.conn.Query(ctx, `
		SELECT id, email, code_hash, expires_at, attempt_count, invalidated_at, created_at
		FROM otp_tokens
		WHERE email = $1
		  AND invalidated_at IS NULL
		  AND expires_at > NOW()
		  AND attempt_count < 5
		ORDER BY created_at DESC
		LIMIT 1
	`, email)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query OTP token: %s", err.Error()))
	}

	otp, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[OTPToken])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewNotFound("no active OTP found for email")
		}
		return nil, errors.New(fmt.Sprintf("failed to scan OTP token: %s", err.Error()))
	}

	return otp, nil
}

func (s *store) incrementOTPAttempts(ctx context.Context, id uuid.UUID) (int, error) {
	var newCount int
	err := s.conn.QueryRow(ctx,
		"UPDATE otp_tokens SET attempt_count = attempt_count + 1 WHERE id = $1 RETURNING attempt_count",
		id,
	).Scan(&newCount)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("failed to increment OTP attempt count: %s", err.Error()))
	}
	return newCount, nil
}

func (s *store) invalidateOTP(ctx context.Context, id uuid.UUID) error {
	_, err := s.conn.Exec(ctx,
		"UPDATE otp_tokens SET invalidated_at = NOW() WHERE id = $1",
		id,
	)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to invalidate OTP token: %s", err.Error()))
	}
	return nil
}

func (s *store) createOTPToken(ctx context.Context, email, codeHash string, expiresAt time.Time) error {
	_, err := s.conn.Exec(ctx,
		"INSERT INTO otp_tokens(email, code_hash, expires_at) VALUES($1, $2, $3)",
		email, codeHash, expiresAt,
	)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create OTP token: %s", err.Error()))
	}
	return nil
}

func (s *store) createUser(ctx context.Context, email string) (*User, error) {
	rows, err := s.conn.Query(
		ctx,
		"INSERT INTO users(email, status) VALUES($1, 'pending') RETURNING *",
		email,
	)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to insert user: %s", err.Error()))
	}

	u, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[User])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to scan user: %s", err.Error()))
	}

	return u, nil
}
