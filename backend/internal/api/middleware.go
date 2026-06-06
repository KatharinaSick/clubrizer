package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/katharinasick/clubrizer/internal/app"
	"github.com/katharinasick/clubrizer/internal/users"
	"net/http"
)

func authenticated(cfg app.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessTokenString := r.Header.Get(cfg.JWT.AccessToken.HeaderName)
		if accessTokenString == "" {
			http.Error(w, "no access token set", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(accessTokenString, &users.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWT.AccessToken.SecretKey), nil
		})

		if err != nil || !token.Valid {
			if errors.Is(err, jwt.ErrTokenExpired) {
				http.Error(w, "access token expired", http.StatusUnauthorized)
			} else {
				http.Error(w, "invalid access token", http.StatusUnauthorized)
			}
			return
		}

		claims, ok := token.Claims.(*users.Claims)
		if !ok {
			http.Error(w, "invalid access token claims", http.StatusBadRequest)
			return
		}

		if claims.Status != "approved" {
			http.Error(w, "account not approved", http.StatusForbidden)
			return
		}

		// Add user claims to context for later use
		ctx := context.WithValue(r.Context(), cfg.JWT.User.Key, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
