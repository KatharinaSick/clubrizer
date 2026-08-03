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

// authenticated gates a route to fully approved accounts — the default for almost
// everything.
func authenticated(cfg app.Config, next http.Handler) http.Handler {
	return withClaims(cfg, func(status string) bool { return status == "approved" }, next)
}

// notRejected gates reading the account's own kids to any non-rejected account:
// 'onboarding' (adding kids during setup), 'pending' (the waiting screen lists the
// submitted kids for reassurance), and 'approved' (managing an established account).
// Only 'rejected' is excluded — a declined account has nothing to view. Kid RSVPs are
// separately gated on the kid being approved.
func notRejected(cfg app.Config, next http.Handler) http.Handler {
	return withClaims(cfg, func(status string) bool {
		return status == "onboarding" || status == "pending" || status == "approved"
	}, next)
}

// onboardingOnly gates routes to accounts that are still in onboarding — used for the
// "replace my kids" endpoint, whose wholesale replace must never be reachable once an
// account is approved (where it could delete real kids and their event responses).
func onboardingOnly(cfg app.Config, next http.Handler) http.Handler {
	return withClaims(cfg, func(status string) bool { return status == "onboarding" }, next)
}

// withClaims validates the access token and, if statusAllowed accepts the account's
// status, injects the claims into the request context. Otherwise it rejects with 403.
func withClaims(cfg app.Config, statusAllowed func(status string) bool, next http.Handler) http.Handler {
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

		if !statusAllowed(claims.Status) {
			http.Error(w, "account not approved", http.StatusForbidden)
			return
		}

		// Add user claims to context for later use
		ctx := context.WithValue(r.Context(), cfg.JWT.User.Key, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
