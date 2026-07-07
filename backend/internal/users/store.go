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
	cutoff := time.Now().Add(-time.Duration(s.cfg.OTP.WindowHours) * time.Hour)
	err := s.conn.QueryRow(ctx,
		"SELECT COUNT(*) FROM otp_tokens WHERE email = $1 AND created_at > $2",
		email, cutoff,
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

// consumeOTP atomically marks an OTP token as used. It reports whether this
// call was the one that consumed it (true), or whether it had already been
// consumed by a concurrent request (false). Callers may proceed in either case:
// reaching this point means the correct code was presented for a token that was
// active when it was read, so a concurrent consume is a duplicate of a
// legitimate login rather than a reason to reject it.
func (s *store) consumeOTP(ctx context.Context, id uuid.UUID) (bool, error) {
	ct, err := s.conn.Exec(ctx,
		"UPDATE otp_tokens SET invalidated_at = NOW() WHERE id = $1 AND invalidated_at IS NULL",
		id,
	)
	if err != nil {
		return false, errors.New(fmt.Sprintf("failed to consume OTP token: %s", err.Error()))
	}
	return ct.RowsAffected() > 0, nil
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

func (s *store) updateUserProfile(ctx context.Context, id uuid.UUID, firstName, lastName string, nickName *string) (*User, error) {
	rows, err := s.conn.Query(ctx,
		"UPDATE users SET given_name = $1, family_name = $2, nick_name = $3 WHERE id = $4 RETURNING *",
		firstName, lastName, nickName, id,
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to update user profile: %s", err.Error()))
	}

	u, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[User])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to scan user: %s", err.Error()))
	}

	return u, nil
}

func (s *store) updateUserPicture(ctx context.Context, id uuid.UUID, pictureURL string) (*User, error) {
	rows, err := s.conn.Query(ctx,
		"UPDATE users SET picture = $1 WHERE id = $2 RETURNING *",
		pictureURL, id,
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to update user picture: %s", err.Error()))
	}

	u, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[User])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to scan user: %s", err.Error()))
	}

	return u, nil
}

func (s *store) getRolesByUserID(ctx context.Context, userID uuid.UUID) ([]*Role, error) {
	rows, err := s.conn.Query(ctx, `
		SELECT r.id, r.name
		FROM user_roles ur
		JOIN roles r ON r.id = ur.role_id
		WHERE ur.user_id = $1 AND r.name != 'member'
		ORDER BY r.name
	`, userID)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query roles: %s", err.Error()))
	}
	defer rows.Close()

	var roles []*Role
	for rows.Next() {
		var r Role
		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			return nil, errors.New(fmt.Sprintf("failed to scan role: %s", err.Error()))
		}
		roles = append(roles, &r)
	}
	if rows.Err() != nil {
		return nil, errors.New(fmt.Sprintf("rows error: %s", rows.Err().Error()))
	}
	return roles, nil
}

func (s *store) createUser(ctx context.Context, email string) (*User, error) {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return nil, errors.New("failed to begin transaction: " + err.Error())
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(
		ctx,
		"INSERT INTO users(email, status) VALUES($1, 'pending') ON CONFLICT (email) DO NOTHING RETURNING *",
		email,
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to insert user: %s", err.Error()))
	}

	u, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// A concurrent verify for the same email already created the user
			// (e.g. a duplicate request during first-time registration). Load
			// the existing user instead of failing — the member role was
			// already granted by the request that won the insert.
			existingRows, selErr := tx.Query(ctx, "SELECT * FROM users WHERE email = $1", email)
			if selErr != nil {
				return nil, errors.New(fmt.Sprintf("failed to query user: %s", selErr.Error()))
			}
			u, selErr = pgx.CollectOneRow(existingRows, pgx.RowToAddrOfStructByName[User])
			if selErr != nil {
				return nil, errors.New(fmt.Sprintf("failed to scan user: %s", selErr.Error()))
			}
			if err := tx.Commit(ctx); err != nil {
				return nil, errors.New(fmt.Sprintf("failed to commit transaction: %s", err.Error()))
			}
			return u, nil
		}
		return nil, errors.New(fmt.Sprintf("failed to scan user: %s", err.Error()))
	}

	_, err = tx.Exec(ctx,
		"INSERT INTO user_roles (user_id, role_id) SELECT $1, id FROM roles WHERE name = 'member'",
		u.ID,
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to grant member role: %s", err.Error()))
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to commit transaction: %s", err.Error()))
	}

	return u, nil
}
