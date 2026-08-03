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
		"INSERT INTO users(email, status) VALUES($1, 'onboarding') ON CONFLICT (email) DO NOTHING RETURNING *",
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
			// the existing user instead of failing.
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

	// No role is granted here: the account is still in onboarding and hasn't chosen its
	// type yet. The member-role grant is settled once when onboarding completes, tracking
	// self_participates (participants are members, guardians are not) — see submitForApproval.

	if err := tx.Commit(ctx); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to commit transaction: %s", err.Error()))
	}

	return u, nil
}

const kidColumns = "id, user_id, given_name, family_name, picture, status, created_at"

// getKidsByUserID returns the account's kids for management, excluding rejected and
// removed ones: a declined kid is treated as non-existent (there is no resubmit flow
// yet), and a removed kid is soft-deleted (deleted_at set) so its past event responses
// survive — both are hidden here rather than lingering in the list.
func (s *store) getKidsByUserID(ctx context.Context, userID uuid.UUID) ([]*Kid, error) {
	rows, err := s.conn.Query(ctx,
		"SELECT "+kidColumns+" FROM kids WHERE user_id = $1 AND status != 'rejected' AND deleted_at IS NULL ORDER BY created_at",
		userID,
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query kids: %s", err.Error()))
	}

	kids, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Kid])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to scan kids: %s", err.Error()))
	}
	return kids, nil
}

func (s *store) createKid(ctx context.Context, userID uuid.UUID, givenName, familyName string) (*Kid, error) {
	rows, err := s.conn.Query(ctx,
		"INSERT INTO kids (user_id, given_name, family_name) VALUES ($1, $2, NULLIF($3, '')) RETURNING "+kidColumns,
		userID, givenName, familyName,
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to insert kid: %s", err.Error()))
	}

	kid, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[Kid])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to scan kid: %s", err.Error()))
	}
	return kid, nil
}

// replaceKids sets the account's kids to exactly the given given-names, in one
// transaction (clearing any existing kids first). Used only during onboarding, where
// every kid is a freshly-typed pending row with no event responses — so a wholesale
// replace is safe and lets the onboarding screen add/rename/remove kids in a single
// request. New kids are created pending; family name and picture are filled in after
// approval.
func (s *store) replaceKids(ctx context.Context, userID uuid.UUID, givenNames []string) ([]*Kid, error) {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return nil, errors.New("failed to begin transaction: " + err.Error())
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, "DELETE FROM kids WHERE user_id = $1", userID); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to clear kids: %s", err.Error()))
	}

	rows, err := tx.Query(ctx,
		"INSERT INTO kids (user_id, given_name) SELECT $1, name FROM unnest($2::text[]) AS name RETURNING "+kidColumns,
		userID, givenNames,
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to insert kids: %s", err.Error()))
	}
	kids, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Kid])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to scan kids: %s", err.Error()))
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to commit transaction: %s", err.Error()))
	}
	return kids, nil
}

// updateKid updates a kid's name, scoped to the owning parent. If no row matches
// both the kid id and the owner it returns NotFound, so a guessed id belonging to
// another parent is indistinguishable from a non-existent one.
func (s *store) updateKid(ctx context.Context, kidID, userID uuid.UUID, givenName, familyName string) (*Kid, error) {
	rows, err := s.conn.Query(ctx,
		"UPDATE kids SET given_name = $1, family_name = NULLIF($2, '') WHERE id = $3 AND user_id = $4 AND deleted_at IS NULL RETURNING "+kidColumns,
		givenName, familyName, kidID, userID,
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to update kid: %s", err.Error()))
	}

	kid, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[Kid])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewNotFound(fmt.Sprintf("kid with id %s not found", kidID))
		}
		return nil, errors.New(fmt.Sprintf("failed to scan kid: %s", err.Error()))
	}
	return kid, nil
}

func (s *store) updateKidPicture(ctx context.Context, kidID, userID uuid.UUID, pictureURL string) (*Kid, error) {
	rows, err := s.conn.Query(ctx,
		"UPDATE kids SET picture = $1 WHERE id = $2 AND user_id = $3 AND deleted_at IS NULL RETURNING "+kidColumns,
		pictureURL, kidID, userID,
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to update kid picture: %s", err.Error()))
	}

	kid, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[Kid])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewNotFound(fmt.Sprintf("kid with id %s not found", kidID))
		}
		return nil, errors.New(fmt.Sprintf("failed to scan kid: %s", err.Error()))
	}
	return kid, nil
}

// deleteKid soft-deletes a kid scoped to the owning parent by stamping deleted_at. The
// row and its event_responses are kept on purpose, so the kid still shows in the attendee
// lists of events they already responded to; they're just hidden from the parent's
// management list and RSVP controls (deleted_at IS NULL filters) and can't be edited or
// receive new responses. Deleting an already-deleted kid returns NotFound, matching the
// old hard-delete behaviour.
func (s *store) deleteKid(ctx context.Context, kidID, userID uuid.UUID) error {
	ct, err := s.conn.Exec(ctx,
		"UPDATE kids SET deleted_at = NOW() WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL",
		kidID, userID,
	)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to delete kid: %s", err.Error()))
	}
	if ct.RowsAffected() == 0 {
		return apperrors.NewNotFound(fmt.Sprintf("kid with id %s not found", kidID))
	}
	return nil
}

// setSelfParticipates records the account's self_participates flag. The member-role
// grant is intentionally NOT touched here: during onboarding the role is irrelevant
// (nothing member-gated is reachable yet), and it is settled once when onboarding
// completes — see submitForApproval. Keeping it to a single point avoids a
// grant-here/revoke-there split that drifts out of sync.
func (s *store) setSelfParticipates(ctx context.Context, userID uuid.UUID, selfParticipates bool) (*User, error) {
	rows, err := s.conn.Query(ctx,
		"UPDATE users SET self_participates = $1 WHERE id = $2 RETURNING *",
		selfParticipates, userID,
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to update self_participates: %s", err.Error()))
	}
	u, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[User])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to scan user: %s", err.Error()))
	}
	return u, nil
}

func (s *store) countKidsByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := s.conn.QueryRow(ctx, "SELECT count(*) FROM kids WHERE user_id = $1 AND deleted_at IS NULL", userID).Scan(&count)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("failed to count kids: %s", err.Error()))
	}
	return count, nil
}

// submitForApproval moves the account from 'onboarding' to 'pending' and settles the
// member-role grant to match the account type, in one transaction. This is the single
// place the grant is decided: participants ("just me" / "me & my kids") become members,
// guardians ("only my kids") do not. Doing it here — rather than granting at account
// creation and revoking later — means a guardian never holds the member role at all.
//
// The grant is driven off the authoritative users.self_participates column (not a value
// passed in from the caller's claims), so a stale or forged token can't hand a guardian
// the member role. Both statements always run; the WHERE clauses make exactly one take
// effect for the account's current type.
func (s *store) submitForApproval(ctx context.Context, userID uuid.UUID) (*User, error) {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return nil, errors.New("failed to begin transaction: " + err.Error())
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, "UPDATE users SET status = 'pending' WHERE id = $1", userID); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to update user status: %s", err.Error()))
	}

	// Grant member iff the account participates itself.
	if _, err := tx.Exec(ctx,
		`INSERT INTO user_roles (user_id, role_id)
		 SELECT u.id, r.id FROM users u JOIN roles r ON r.name = 'member'
		 WHERE u.id = $1 AND u.self_participates
		 ON CONFLICT (user_id, role_id) DO NOTHING`,
		userID,
	); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to grant member role: %s", err.Error()))
	}
	// Revoke member iff the account is a guardian (belt-and-braces: guardians never
	// held it, but this keeps the invariant true regardless of prior state).
	if _, err := tx.Exec(ctx,
		`DELETE FROM user_roles
		 WHERE user_id = $1
		   AND role_id = (SELECT id FROM roles WHERE name = 'member')
		   AND NOT (SELECT self_participates FROM users WHERE id = $1)`,
		userID,
	); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to revoke member role: %s", err.Error()))
	}

	rows, err := tx.Query(ctx, "SELECT * FROM users WHERE id = $1", userID)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query user: %s", err.Error()))
	}
	u, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[User])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to scan user: %s", err.Error()))
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to commit transaction: %s", err.Error()))
	}
	return u, nil
}

// getOwnedKid returns the kid only if it belongs to userID and hasn't been removed (any
// status). Callers that also require approval (e.g. RSVP) check Status themselves.
func (s *store) getOwnedKid(ctx context.Context, kidID, userID uuid.UUID) (*Kid, error) {
	rows, err := s.conn.Query(ctx,
		"SELECT "+kidColumns+" FROM kids WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL",
		kidID, userID,
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query kid: %s", err.Error()))
	}
	kid, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[Kid])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewNotFound(fmt.Sprintf("kid with id %s not found", kidID))
		}
		return nil, errors.New(fmt.Sprintf("failed to scan kid: %s", err.Error()))
	}
	return kid, nil
}
