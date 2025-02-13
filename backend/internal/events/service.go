package events

import (
	"github.com/katharinasick/clubrizer/internal/app"
)

// Service is responsible for everything related to handling events.
type Service struct {
	log app.Logger
	cfg *app.Config
}

func NewService(log app.Logger, cfg *app.Config) *Service {
	return &Service{log, cfg}
}
