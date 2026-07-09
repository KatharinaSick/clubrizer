package api

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/katharinasick/clubrizer/internal/app"
	"github.com/katharinasick/clubrizer/internal/events"
	"github.com/katharinasick/clubrizer/internal/users"
	"github.com/rs/cors"
)

func NewHandler(
	log app.Logger,
	cfg app.Config,
	userService userService,
	eventsService eventsService,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, cfg, userService, eventsService)

	handler := cors.New(cors.Options{
		AllowedOrigins: cfg.Cors.AllowedOrigins,
		AllowedHeaders: cfg.Cors.AllowedHeaders,
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowCredentials: true,
	}).Handler(mux)

	return requestLogger(log, handler)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
	errMsg strings.Builder
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status >= 400 {
		r.errMsg.Write(b)
	}
	return r.ResponseWriter.Write(b)
}

func requestLogger(log app.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		args := []any{
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
		}
		if rec.status >= 400 {
			args = append(args, "error", strings.TrimSpace(rec.errMsg.String()))
		}
		if rec.status >= 500 {
			log.Error("request", args...)
		} else if rec.status >= 400 {
			log.Warn("request", args...)
		} else {
			log.Info("request", args...)
		}
	})
}

func addRoutes(
	mux *http.ServeMux,
	cfg app.Config,
	userService userService,
	eventsService eventsService,
) {
	// Authentication & Users
	mux.Handle("POST /auth/otp", handleWithBody(userService.RequestOTP))
	mux.Handle("POST /auth/verify", handleWithBodyAndReturnRefreshToken(cfg, userService.VerifyOTP))
	mux.Handle("POST /auth/refresh", handleWithRefreshToken(cfg, userService.RefreshTokens))
	mux.Handle("POST /auth/logout", handleLogout(cfg))

	// Users
	mux.Handle("GET /users/me/roles", authenticated(cfg, handleAndReturnList(userService.GetMyRoles)))
	mux.Handle("PATCH /users/me/profile", authenticated(cfg, handleWithBodyAndReturnRefreshToken(cfg, userService.UpdateProfile)))
	mux.Handle("POST /users/me/picture", authenticated(cfg, handleProfilePicture(cfg, userService.UpdateProfilePicture)))

	// Events
	mux.Handle("GET /events/categories", authenticated(cfg, handleAndReturnList(eventsService.ListCategories)))

	mux.Handle("GET /events", authenticated(cfg, handleAndReturnList(eventsService.ListEvents)))
	mux.Handle("GET /events/{id}", authenticated(cfg, handleWithIdAndReturnValue(eventsService.GetEvent)))
	mux.Handle("POST /events", authenticated(cfg, handleWithBodyAndReturnValue(eventsService.CreateEvent)))
	mux.Handle("PUT /events/{id}/response", authenticated(cfg, handleWithIdAndBody(eventsService.UpsertEventResponse)))
	mux.Handle("DELETE /events/{id}", authenticated(cfg, handleWithId(eventsService.DeleteEvent)))
}

type userService interface {
	RequestOTP(ctx context.Context, req users.RequestOTPRequest) error
	VerifyOTP(ctx context.Context, req users.VerifyOTPRequest) (*users.VerifyOTPResponse, *users.RefreshTokenInfo, error)
	GetMyRoles(ctx context.Context) ([]*users.Role, error)
	UpdateProfile(ctx context.Context, req users.UpdateProfileRequest) (*users.UpdateProfileResponse, *users.RefreshTokenInfo, error)
	UpdateProfilePicture(ctx context.Context, contentType string, data io.Reader) (*users.UpdateProfilePictureResponse, *users.RefreshTokenInfo, error)
	RefreshTokens(ctx context.Context, t users.RefreshTokenInfo) (*users.RefreshTokensResponse, *users.RefreshTokenInfo, error)
}

type eventsService interface {
	ListCategories(ctx context.Context) ([]*events.Category, error)
	ListEvents(ctx context.Context) ([]*events.Event, error)
	GetEvent(ctx context.Context, id string) (*events.Event, error)
	CreateEvent(ctx context.Context, e events.Event) (*events.Event, error)
	UpsertEventResponse(ctx context.Context, eventId string, req events.UpsertEventResponseRequest) error
	DeleteEvent(ctx context.Context, id string) error
}
