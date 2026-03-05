package events

import (
	"context"
	"github.com/katharinasick/clubrizer/internal/apperrors"
	"time"
)

func (s *Service) ListEvents(ctx context.Context) ([]*Event, error) {
	events, err := s.store.getAllEvents(ctx)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *Service) CreateEvent(ctx context.Context, e Event) (*Event, error) {
	if e.StartTime.Before(time.Now()) {
		return nil, apperrors.NewBadRequest("start time must be in the future", nil)
	}

	id, err := s.store.createEvent(ctx, &e)
	if err != nil {
		return nil, err
	}
	e.ID = id
	return &e, nil
}
