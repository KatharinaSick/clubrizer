package events

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/katharinasick/clubrizer/internal/apperrors"
	"github.com/katharinasick/clubrizer/internal/rbac"
	"github.com/katharinasick/clubrizer/internal/users"
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

	canDelete, err := s.canDeleteEvent(ctx, userId, event.CreatedBy)
	if err != nil {
		return nil, err
	}
	event.CanDelete = &canDelete

	return event, nil
}

// canDeleteEvent reports whether the user may delete the event with the given creator.
// The event's owner can always delete it; everyone else needs the delete-any permission
// (which admins bypass).
func (s *Service) canDeleteEvent(ctx context.Context, userId, creatorId uuid.UUID) (bool, error) {
	if userId == creatorId {
		return true, nil
	}
	return s.rbac.IsAuthorized(ctx, userId, rbac.PermissionEventsDeleteAny)
}

func (s *Service) DeleteEvent(ctx context.Context, id string) error {
	uuidId, err := uuid.Parse(id)
	if err != nil {
		return apperrors.NewBadRequest("invalid event id", err)
	}

	event, err := s.store.getEventById(ctx, uuidId)
	if err != nil {
		return err
	}

	userId := ctx.Value(s.cfg.JWT.User.Key).(*users.Claims).ID

	authorized, err := s.canDeleteEvent(ctx, userId, event.CreatedBy)
	if err != nil {
		return err
	}
	if !authorized {
		return apperrors.NewForbidden("you are not allowed to delete this event")
	}

	return s.store.deleteEvent(ctx, uuidId)
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

	userId := ctx.Value(s.cfg.JWT.User.Key).(*users.Claims).ID

	authorized, err := s.rbac.IsAuthorizedToCreateEvent(ctx, userId, e.CategoryID)
	if err != nil {
		return nil, err
	}
	if !authorized {
		return nil, apperrors.NewForbidden("you are not allowed to create events in this category")
	}

	id, err := s.store.createEvent(ctx, &e)
	if err != nil {
		return nil, err
	}
	e.ID = id

	if err := s.store.upsertEventResponse(ctx, id, userId, true); err != nil {
		return nil, err
	}

	return &e, nil
}
