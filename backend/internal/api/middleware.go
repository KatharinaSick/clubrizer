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

// requirePermission gates a route to approved accounts that additionally hold the given
// permission (admins bypass the permission check inside the authorizer). It layers on top
// of authenticated, so the status gate and claims injection are shared and only the
// permission varies per route. Use it for pure permission gates; checks that depend on a
// specific resource (e.g. an event's owner) or on the request body still belong in the
// service, which cannot be expressed at route level.
func requirePermission(cfg app.Config, auth authorizer, permission string, next http.Handler) http.Handler {
	return authenticated(cfg, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(cfg.JWT.User.Key).(*users.Claims)

		authorized, err := auth.IsAuthorized(r.Context(), claims.ID, permission)
		if err != nil {
			http.Error(w, "failed to check permission", http.StatusInternalServerError)
			return
		}
		if !authorized {
			http.Error(w, "insufficient permissions", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}))
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
