package users

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
}
