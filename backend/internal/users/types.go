package users

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Status is the approval lifecycle state of an account or kid. A new signup starts at
// StatusOnboarding, becomes StatusPending when submitted for approval, and is then moved
// by an admin to StatusApproved or StatusRejected. Stored as a plain text/enum column in
// the DB; this named type keeps the values from being loose string literals in Go code.
type Status string

const (
	StatusOnboarding Status = "onboarding"
	StatusPending    Status = "pending"
	StatusApproved   Status = "approved"
	StatusRejected   Status = "rejected"
)

type User struct {
	ID               uuid.UUID `db:"id" json:"-"`
	Email            string    `db:"email" json:"email"`
	FamilyName       *string   `db:"family_name" json:"familyName"`
	GivenName        *string   `db:"given_name" json:"givenName"`
	NickName         *string   `db:"nick_name" json:"nickName"`
	Picture          *string   `db:"picture" json:"picture"`
	Status           string    `db:"status" json:"status"`
	SelfParticipates bool      `db:"self_participates" json:"selfParticipates"`
}

// Kid is a participant managed by a parent account. Kids never authenticate — they
// have no email, no roles and no auth of their own. The owning parent is UserID.
type Kid struct {
	ID         uuid.UUID `db:"id" json:"id"`
	UserID     uuid.UUID `db:"user_id" json:"-"`
	GivenName  *string   `db:"given_name" json:"givenName"`
	FamilyName *string   `db:"family_name" json:"familyName"`
	Picture    *string   `db:"picture" json:"picture"`
	Status     string    `db:"status" json:"status"`
	CreatedAt  time.Time `db:"created_at" json:"-"`
}

type OTPToken struct {
	ID            uuid.UUID  `db:"id"`
	Email         string     `db:"email"`
	CodeHash      string     `db:"code_hash"`
	ExpiresAt     time.Time  `db:"expires_at"`
	AttemptCount  int        `db:"attempt_count"`
	InvalidatedAt *time.Time `db:"invalidated_at"`
	CreatedAt     time.Time  `db:"created_at"`
}

type Role struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type Claims struct {
	jwt.RegisteredClaims
	ID               uuid.UUID `json:"id"`
	Email            string    `json:"email"`
	FamilyName       *string   `json:"familyName"`
	GivenName        *string   `json:"givenName"`
	NickName         *string   `json:"nickName"`
	Picture          *string   `json:"picture"`
	Status           string    `json:"status"`
	SelfParticipates bool      `json:"selfParticipates"`
	// Permissions are the caller's effective permission keys, included so the frontend can
	// show or hide UI without an extra request. A hint only — backend authorization always
	// re-checks against the DB, since a token can be up to its lifetime out of date.
	Permissions []string `json:"permissions"`
}
