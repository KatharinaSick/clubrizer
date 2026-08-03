package users

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
//
// The transition is guarded on the account still being in 'onboarding', making the whole
// operation idempotent: a double-submit or a replay with a still-valid onboarding token
// flips no row the second time. It returns transitioned=false in that case so the caller
// skips re-notifying admins (and no role work is redone).
func (s *store) submitForApproval(ctx context.Context, userID uuid.UUID) (*User, bool, error) {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return nil, false, errors.New("failed to begin transaction: " + err.Error())
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, "UPDATE users SET status = 'pending' WHERE id = $1 AND status = 'onboarding'", userID)
	if err != nil {
		return nil, false, errors.New(fmt.Sprintf("failed to update user status: %s", err.Error()))
	}
	transitioned := tag.RowsAffected() > 0

	if transitioned {
		// Grant member iff the account participates itself.
		if _, err := tx.Exec(ctx,
			`INSERT INTO user_roles (user_id, role_id)
			 SELECT u.id, r.id FROM users u JOIN roles r ON r.name = 'member'
			 WHERE u.id = $1 AND u.self_participates
			 ON CONFLICT (user_id, role_id) DO NOTHING`,
			userID,
		); err != nil {
			return nil, false, errors.New(fmt.Sprintf("failed to grant member role: %s", err.Error()))
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
			return nil, false, errors.New(fmt.Sprintf("failed to revoke member role: %s", err.Error()))
		}
	}

	rows, err := tx.Query(ctx, "SELECT * FROM users WHERE id = $1", userID)
	if err != nil {
		return nil, false, errors.New(fmt.Sprintf("failed to query user: %s", err.Error()))
	}
	u, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[User])
	if err != nil {
		return nil, false, errors.New(fmt.Sprintf("failed to scan user: %s", err.Error()))
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, false, errors.New(fmt.Sprintf("failed to commit transaction: %s", err.Error()))
	}
	return u, transitioned, nil
}

// listApprovals returns the approval queue as one row per account: every account still
// awaiting approval ('pending'), plus every approved account that has newly-added kids
// awaiting approval. Onboarding accounts are excluded — they haven't submitted yet, even
// though their kids are already 'pending' — as are rejected accounts. Each account's
// pending kids are attached in a second query.
func (s *store) listApprovals(ctx context.Context) ([]*ApprovalRequest, error) {
	rows, err := s.conn.Query(ctx, `
		SELECT id, email, given_name, family_name, nick_name, picture, status, self_participates
		FROM users u
		WHERE u.status = 'pending'
		   OR (u.status = 'approved' AND EXISTS (
		        SELECT 1 FROM kids k
		        WHERE k.user_id = u.id AND k.status = 'pending' AND k.deleted_at IS NULL))
		ORDER BY email
	`)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query approvals: %s", err.Error()))
	}
	defer rows.Close()

	var requests []*ApprovalRequest
	byID := make(map[uuid.UUID]*ApprovalRequest)
	for rows.Next() {
		f := &ApprovalRequest{PendingKids: []*ApprovalKid{}}
		if err := rows.Scan(&f.UserID, &f.Email, &f.GivenName, &f.FamilyName, &f.NickName, &f.Picture, &f.Status, &f.SelfParticipates); err != nil {
			return nil, errors.New(fmt.Sprintf("failed to scan approval: %s", err.Error()))
		}
		requests = append(requests, f)
		byID[f.UserID] = f
	}
	if rows.Err() != nil {
		return nil, errors.New(fmt.Sprintf("rows error: %s", rows.Err().Error()))
	}

	if len(requests) == 0 {
		return requests, nil
	}

	ids := make([]uuid.UUID, 0, len(byID))
	for id := range byID {
		ids = append(ids, id)
	}

	kidRows, err := s.conn.Query(ctx, `
		SELECT id, user_id, given_name, family_name, picture
		FROM kids
		WHERE user_id = ANY($1) AND status = 'pending' AND deleted_at IS NULL
		ORDER BY created_at
	`, ids)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query pending kids: %s", err.Error()))
	}
	defer kidRows.Close()

	for kidRows.Next() {
		var userID uuid.UUID
		k := &ApprovalKid{}
		if err := kidRows.Scan(&k.ID, &userID, &k.GivenName, &k.FamilyName, &k.Picture); err != nil {
			return nil, errors.New(fmt.Sprintf("failed to scan pending kid: %s", err.Error()))
		}
		if f, ok := byID[userID]; ok {
			f.PendingKids = append(f.PendingKids, k)
		}
	}
	if kidRows.Err() != nil {
		return nil, errors.New(fmt.Sprintf("rows error: %s", kidRows.Err().Error()))
	}

	return requests, nil
}

// decideApprovals sets the given accounts and kids to decision (StatusApproved or
// StatusRejected) in one all-or-nothing transaction, returning who to notify (see
// decisionRecipients). Only rows currently 'pending' are eligible;
// if any requested id is not pending (stale, already decided, bogus, or — for kids —
// belonging to an account that isn't approved) the affected-row count won't match and the
// whole batch is rolled back. When approving, a kid is only eligible if its parent is
// already approved or is being approved in this same batch — the users UPDATE runs first, so
// batched parents already read as 'approved' within the transaction. This guarantees an
// approved kid never hangs under a non-approved parent. Rejecting a kid carries no such
// parent guard, and rejecting an account additionally cascades to its still-pending kids so
// none are stranded (see below).
func (s *store) decideApprovals(ctx context.Context, userIDs, kidIDs []uuid.UUID, decision Status) (decisionRecipients, error) {
	// De-duplicate first: a repeated id would otherwise undercount against len() below and
	// trip the all-or-nothing check on input that is actually valid.
	userIDs = dedupeUUIDs(userIDs)
	kidIDs = dedupeUUIDs(kidIDs)

	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return decisionRecipients{}, errors.New("failed to begin transaction: " + err.Error())
	}
	defer tx.Rollback(ctx)

	var accountEmails []string
	if len(userIDs) > 0 {
		rows, err := tx.Query(ctx,
			"UPDATE users SET status = $1 WHERE id = ANY($2) AND status = 'pending' RETURNING email",
			string(decision), userIDs,
		)
		if err != nil {
			return decisionRecipients{}, errors.New(fmt.Sprintf("failed to update users: %s", err.Error()))
		}
		accountEmails, err = pgx.CollectRows(rows, pgx.RowTo[string])
		if err != nil {
			return decisionRecipients{}, errors.New(fmt.Sprintf("failed to collect updated emails: %s", err.Error()))
		}
		if len(accountEmails) != len(userIDs) {
			return decisionRecipients{}, apperrors.NewBadRequest("one or more accounts are not awaiting approval", nil)
		}
	}

	if len(kidIDs) > 0 {
		var tag pgconn.CommandTag
		if decision == StatusApproved {
			tag, err = tx.Exec(ctx, `
				UPDATE kids SET status = 'approved'
				WHERE id = ANY($1) AND status = 'pending' AND deleted_at IS NULL
				  AND user_id IN (SELECT id FROM users WHERE status = 'approved')
			`, kidIDs)
		} else {
			tag, err = tx.Exec(ctx, `
				UPDATE kids SET status = 'rejected'
				WHERE id = ANY($1) AND status = 'pending' AND deleted_at IS NULL
			`, kidIDs)
		}
		if err != nil {
			return decisionRecipients{}, errors.New(fmt.Sprintf("failed to update kids: %s", err.Error()))
		}
		if int(tag.RowsAffected()) != len(kidIDs) {
			return decisionRecipients{}, apperrors.NewBadRequest("one or more kids cannot be decided (not awaiting approval, or their account is not approved)", nil)
		}
	}

	// Rejecting an account cascades to its still-pending kids. Otherwise a kid left pending
	// under a rejected parent is stranded: listApprovals surfaces neither rejected accounts
	// nor their kids, and the approve path requires an approved parent — so the kid would be
	// invisible in the queue and impossible to ever decide. Runs after the explicit kid
	// updates above and is deliberately not counted against kidIDs; it only mops up pending
	// kids not already handled in this batch.
	if decision == StatusRejected && len(userIDs) > 0 {
		if _, err := tx.Exec(ctx, `
			UPDATE kids SET status = 'rejected'
			WHERE user_id = ANY($1) AND status = 'pending' AND deleted_at IS NULL
		`, userIDs); err != nil {
			return decisionRecipients{}, errors.New(fmt.Sprintf("failed to reject kids of rejected accounts: %s", err.Error()))
		}
	}

	// Parents to notify about an individual kid decision: the decided kids whose account is
	// NOT in this batch (a later-added kid on an already-approved account). Kids decided
	// alongside their own account — a whole-family batch — are excluded here, because the
	// single account-level email already covers the family.
	kidParents, err := collectKidParentNotifications(ctx, tx, kidIDs, userIDs)
	if err != nil {
		return decisionRecipients{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return decisionRecipients{}, errors.New(fmt.Sprintf("failed to commit transaction: %s", err.Error()))
	}
	return decisionRecipients{accountEmails: accountEmails, kidParents: kidParents}, nil
}

// collectKidParentNotifications returns the parent email + kid name for each decided kid
// (in kidIDs) whose account is not in userIDs. Runs inside the decision transaction so it
// closes its rows before the caller commits.
func collectKidParentNotifications(ctx context.Context, tx pgx.Tx, kidIDs, userIDs []uuid.UUID) ([]kidParentNotification, error) {
	if len(kidIDs) == 0 {
		return nil, nil
	}

	// A nil userIDs encodes such that `k.user_id <> ALL($2)` evaluates to NULL, excluding
	// every row — which would silently skip all parent emails. Normalize to an empty array
	// so a kid-only batch (parent already approved, not in the batch) still notifies parents.
	if userIDs == nil {
		userIDs = []uuid.UUID{}
	}

	rows, err := tx.Query(ctx, `
		SELECT u.email, COALESCE(k.given_name, '')
		FROM kids k
		JOIN users u ON u.id = k.user_id
		WHERE k.id = ANY($1) AND k.user_id <> ALL($2)
	`, kidIDs, userIDs)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query kid parents: %s", err.Error()))
	}
	defer rows.Close()

	var out []kidParentNotification
	for rows.Next() {
		var n kidParentNotification
		if err := rows.Scan(&n.parentEmail, &n.kidName); err != nil {
			return nil, errors.New(fmt.Sprintf("failed to scan kid parent: %s", err.Error()))
		}
		out = append(out, n)
	}
	if rows.Err() != nil {
		return nil, errors.New(fmt.Sprintf("rows error: %s", rows.Err().Error()))
	}
	return out, nil
}

// dedupeUUIDs returns ids with duplicates removed, preserving first-seen order.
func dedupeUUIDs(ids []uuid.UUID) []uuid.UUID {
	if len(ids) < 2 {
		return ids
	}
	seen := make(map[uuid.UUID]struct{}, len(ids))
	out := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

// getAdminEmails returns the emails of all approved admins, used to notify them of new
// approval requests.
func (s *store) getAdminEmails(ctx context.Context) ([]string, error) {
	rows, err := s.conn.Query(ctx, `
		SELECT u.email
		FROM users u
		JOIN user_roles ur ON ur.user_id = u.id
		JOIN roles r ON r.id = ur.role_id
		WHERE r.name = 'admin' AND u.status = 'approved'
	`)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query admin emails: %s", err.Error()))
	}
	emails, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to collect admin emails: %s", err.Error()))
	}
	return emails, nil
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
