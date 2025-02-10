package users

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/katharinasick/clubrizer/internal/apperrors"
	"github.com/katharinasick/clubrizer/internal/setup"
)

type store struct {
	log  setup.Logger
	cfg  *setup.Config
	conn *pgxpool.Pool
}

func newStore(log setup.Logger, cfg *setup.Config, conn *pgxpool.Pool) *store {
	return &store{log, cfg, conn}
}

func (s *store) getUserByMail(ctx context.Context, email string) (*User, error) {
	rows, err := s.conn.Query(ctx, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to query user: %s", err.Error()))
	}

	u, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NewNotFound("user with email not found")
		}
		return nil, errors.New(fmt.Sprintf("failed to scan user: %s", err.Error()))
	}

	return u, nil
}

func (s *store) createUser(ctx context.Context, c *googleTokenClaims) (*User, error) {
	rows, err := s.conn.Query(
		ctx,
		"INSERT INTO users(email, family_name, given_name, nick_name, picture, issuer) VALUES($1, $2, $3, $4, $5, $6) RETURNING *",
		c.Email,
		c.FamilyName,
		c.GivenName,
		c.GivenName,
		c.Picture,
		c.Issuer,
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
