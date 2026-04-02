package users

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `db:"id" json:"-"`
	Email      string    `db:"email" json:"email"`
	FamilyName *string   `db:"family_name" json:"familyName"`
	GivenName  *string   `db:"given_name" json:"givenName"`
	NickName   *string   `db:"nick_name" json:"nickName"`
	Picture    *string   `db:"picture" json:"picture"`
	Status     string    `db:"status" json:"status"`
}

type Claims struct {
	jwt.RegisteredClaims
	ID         uuid.UUID `json:"id"`
	IsNew      bool      `json:"isNew"`
	Email      string    `json:"email"`
	FamilyName *string   `json:"familyName"`
	GivenName  *string   `json:"givenName"`
	NickName   *string   `json:"nickName"`
	Picture    *string   `json:"picture"`
	Status     string    `json:"status"`
}
