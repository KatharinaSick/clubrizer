package api

import (
	"context"
	"github.com/katharinasick/clubrizer/internal/app"
	"github.com/katharinasick/clubrizer/internal/events"
	"github.com/katharinasick/clubrizer/internal/users"
	"github.com/rs/cors"
	"net/http"
)

func NewHandler(
	cfg app.Config,
	userService userService,
	eventsService eventsService,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		cfg,
		userService,
		eventsService,
	)

	handler := cors.
		New(cors.Options{
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
		}).
		Handler(mux)
	return handler
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
	mux.Handle("POST /signin", handleWithBodyAndReturnRefreshToken(cfg, userService.SignInOrRegister))
	mux.Handle("POST /oauth/tokens", handleWithRefreshToken(cfg, userService.RefreshTokens))
	mux.Handle("POST /logout", handleLogout(cfg))

	// Events
	mux.Handle("GET /events/categories", authenticated(cfg, handleAndReturnList(eventsService.ListCategories)))

	mux.Handle("GET /events", authenticated(cfg, handleAndReturnList(eventsService.ListEvents)))
	mux.Handle("GET /events/{id}", authenticated(cfg, handleWithIdAndReturnValue(eventsService.GetEvent)))
	mux.Handle("POST /events", authenticated(cfg, handleWithBodyAndReturnValue(eventsService.CreateEvent)))
	mux.Handle("PUT /events/{id}/response", authenticated(cfg, handleWithIdAndBody(eventsService.UpsertEventResponse)))
}

type userService interface {
	RequestOTP(ctx context.Context, req users.RequestOTPRequest) error
	VerifyOTP(ctx context.Context, req users.VerifyOTPRequest) (*users.VerifyOTPResponse, *users.RefreshTokenInfo, error)
	SignInOrRegister(ctx context.Context, r users.SignInOrRegisterRequest) (*users.SignInOrRegisterResponse, *users.RefreshTokenInfo, error)
	RefreshTokens(ctx context.Context, t users.RefreshTokenInfo) (*users.RefreshTokensResponse, *users.RefreshTokenInfo, error)
}

type eventsService interface {
	ListCategories(ctx context.Context) ([]*events.Category, error)
	ListEvents(ctx context.Context) ([]*events.Event, error)
	GetEvent(ctx context.Context, id string) (*events.Event, error)
	CreateEvent(ctx context.Context, e events.Event) (*events.Event, error)
	UpsertEventResponse(ctx context.Context, eventId string, req events.UpsertEventResponseRequest) error
}
