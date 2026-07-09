package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/katharinasick/clubrizer/internal/app"
	"github.com/katharinasick/clubrizer/internal/apperrors"
	"github.com/katharinasick/clubrizer/internal/users"
)

type handler[Out any] func(context.Context) (*Out, error)
type handlerWithInput[In any] func(context.Context, In) error
type handlerWithInputAndReturnValue[In any, Out any] func(context.Context, In) (*Out, error)
type handlerWithIdAndReturnValue[Out any] func(context.Context, string) (*Out, error)
type handlerWithId func(context.Context, string) error

type handlerWithListReturn[Out any] func(context.Context) ([]*Out, error)

type handlerWithIdAndBody[In any] func(context.Context, string, In) error

var validate = validator.New()

type handlerWithInputAndRefreshTokenReturn[In any, Out any] func(context.Context, In) (*Out, *users.RefreshTokenInfo, error)
type handlerWithRefreshToken[Out any] func(context.Context, users.RefreshTokenInfo) (*Out, *users.RefreshTokenInfo, error)

func handle[Out any](f handler[Out]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		out, err := f(r.Context())
		if err != nil {
			http.Error(w, err.Error(), apperrors.HttpStatusCode(err))
			return
		}

		writeResponse(w, out)
	})
}

func handleAndReturnList[Out any](f handlerWithListReturn[Out]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		out, err := f(r.Context())
		if err != nil {
			http.Error(w, err.Error(), apperrors.HttpStatusCode(err))
			return
		}

		writeResponse(w, out)
	})
}

func handleWithBody[In any](f handlerWithInput[In]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in, err := getAndValidatePayload[In](r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = f(r.Context(), *in)
		if err != nil {
			http.Error(w, err.Error(), apperrors.HttpStatusCode(err))
			return
		}

		ok(w)
	})
}

func handleWithBodyAndReturnValue[In any, Out any](f handlerWithInputAndReturnValue[In, Out]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in, err := getAndValidatePayload[In](r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		out, err := f(r.Context(), *in)
		if err != nil {
			http.Error(w, err.Error(), apperrors.HttpStatusCode(err))
			return
		}

		writeResponse(w, out)
	})
}

func handleWithIdAndReturnValue[Out any](f handlerWithIdAndReturnValue[Out]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		out, err := f(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), apperrors.HttpStatusCode(err))
			return
		}

		writeResponse(w, out)
	})
}

func handleWithId(f handlerWithId) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		if err := f(r.Context(), id); err != nil {
			http.Error(w, err.Error(), apperrors.HttpStatusCode(err))
			return
		}

		ok(w)
	})
}

func handleWithIdAndBody[In any](f handlerWithIdAndBody[In]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		in, err := getAndValidatePayload[In](r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = f(r.Context(), id, *in)
		if err != nil {
			http.Error(w, err.Error(), apperrors.HttpStatusCode(err))
			return
		}

		ok(w)
	})
}

func handleWithBodyAndReturnRefreshToken[In any, Out any](cfg app.Config, f handlerWithInputAndRefreshTokenReturn[In, Out]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in, err := getAndValidatePayload[In](r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		out, rt, err := f(r.Context(), *in)
		if err != nil {
			http.Error(w, err.Error(), apperrors.HttpStatusCode(err))
			return
		}

		setRefreshTokenCookie(w, cfg, rt)
		writeResponse(w, out)
	})
}

func handleWithRefreshToken[Out any](cfg app.Config, f handlerWithRefreshToken[Out]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(cfg.JWT.RefreshToken.CookieName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		out, rt, err := f(r.Context(), users.RefreshTokenInfo{Token: c.Value, Expires: c.Expires})
		if err != nil {
			http.Error(w, err.Error(), apperrors.HttpStatusCode(err))
			return
		}

		setRefreshTokenCookie(w, cfg, rt)
		writeResponse(w, out)
	})
}

func handleProfilePicture(cfg app.Config, f func(context.Context, string, io.Reader) (*users.UpdateProfilePictureResponse, *users.RefreshTokenInfo, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(5 << 20); err != nil { // 5 MB max
			http.Error(w, "file too large or invalid form", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("picture")
		if err != nil {
			http.Error(w, "missing picture field", http.StatusBadRequest)
			return
		}
		defer file.Close()

		contentType := header.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/") {
			http.Error(w, "only image files are accepted", http.StatusBadRequest)
			return
		}

		out, rt, err := f(r.Context(), contentType, file)
		if err != nil {
			http.Error(w, err.Error(), apperrors.HttpStatusCode(err))
			return
		}

		setRefreshTokenCookie(w, cfg, rt)
		writeResponse(w, out)
	})
}

func handleLogout(cfg app.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Clear the refresh token cookie
		setRefreshTokenCookie(w, cfg, &users.RefreshTokenInfo{
			Token:   "",
			Expires: time.Unix(0, 0), // Expire immediately
		})

		ok(w)
	})
}

func setRefreshTokenCookie(w http.ResponseWriter, cfg app.Config, t *users.RefreshTokenInfo) {
	sameSite := http.SameSiteStrictMode
	if cfg.JWT.RefreshToken.SameSite == "Lax" {
		sameSite = http.SameSiteLaxMode
	} else if cfg.JWT.RefreshToken.SameSite == "None" {
		sameSite = http.SameSiteNoneMode
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cfg.JWT.RefreshToken.CookieName,
		Value:    t.Token,
		Expires:  t.Expires,
		Path:     "/",
		HttpOnly: true,
		Secure:   cfg.JWT.RefreshToken.Secure,
		SameSite: sameSite,
	})
}

func getAndValidatePayload[In any](r *http.Request) (*In, error) {
	var in In
	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse json: %s", err))
	}

	if err := validate.Struct(in); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to validate payload: %s", err))
	}

	return &in, nil
}

func writeResponse[Out any](w http.ResponseWriter, out Out) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(out); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

func ok(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
