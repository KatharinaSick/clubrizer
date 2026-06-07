package events

import (
	"context"
	"github.com/google/uuid"
	"github.com/katharinasick/clubrizer/internal/apperrors"
	"github.com/katharinasick/clubrizer/internal/users"
	"time"
)

func (s *Service) ListEvents(ctx context.Context) ([]*Event, error) {
	events, err := s.store.getFutureEvents(ctx)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *Service) GetEvent(ctx context.Context, id string) (*Event, error) {
	uuidId, err := uuid.Parse(id)
	if err != nil {
		return nil, apperrors.NewBadRequest("invalid event id", err)
	}

	event, err := s.store.getEventById(ctx, uuidId)
	if err != nil {
		return nil, err
	}

	userId := ctx.Value(s.cfg.JWT.User.Key).(*users.Claims).ID
	responses, err := s.store.getEventResponses(ctx, uuidId, userId)
	if err != nil {
		return nil, err
	}
	event.Responses = responses

	return event, nil
}

func (s *Service) UpsertEventResponse(ctx context.Context, eventId string, req UpsertEventResponseRequest) error {
	uuidId, err := uuid.Parse(eventId)
	if err != nil {
		return apperrors.NewBadRequest("invalid event id", err)
	}

	event, err := s.store.getEventById(ctx, uuidId)
	if err != nil {
		return err
	}

	if !event.StartTime.After(time.Now()) {
		return apperrors.NewBadRequest("this event has already taken place", nil)
	}

	userId := ctx.Value(s.cfg.JWT.User.Key).(*users.Claims).ID
	return s.store.upsertEventResponse(ctx, uuidId, userId, *req.Response)
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

	userId := ctx.Value(s.cfg.JWT.User.Key).(*users.Claims).ID
	if err := s.store.upsertEventResponse(ctx, id, userId, true); err != nil {
		return nil, err
	}

	return &e, nil
}
