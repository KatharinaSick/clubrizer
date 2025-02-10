package api

import (
	"context"
	"github.com/katharinasick/clubrizer/internal/users"
	"github.com/rs/cors"
	"net/http"
)

func NewHandler(
	userService userService,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		userService,
	)

	handler := cors.Default().Handler(mux)
	return handler
}

func addRoutes(
	mux *http.ServeMux,
	userService userService,
) {
	mux.Handle("POST /users", handleWithBody(userService.SignInOrRegister))
}

type userService interface {
	SignInOrRegister(ctx context.Context, r users.SignInOrRegisterRequest) (*users.SignInOrRegisterResponse, error)
}
