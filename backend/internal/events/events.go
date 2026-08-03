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

	claims := ctx.Value(s.cfg.JWT.User.Key).(*users.Claims)
	userId := claims.ID
	responses, err := s.store.getEventResponses(ctx, uuidId, responder{
		UserID:           claims.ID,
		Name:             displayName(claims),
		Picture:          claims.Picture,
		SelfParticipates: claims.SelfParticipates,
	})
	if err != nil {
		return nil, err
	}
	event.Responses = responses

	canDelete, err := s.canDeleteEvent(ctx, userId, event.CreatedBy)
	if err != nil {
		return nil, err
	}
	event.CanDelete = &canDelete

	canCancel, err := s.canCancelEvent(ctx, userId, event.CreatedBy)
	if err != nil {
		return nil, err
	}
	event.CanCancel = &canCancel

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

// canCancelEvent reports whether the user may cancel the event with the given creator.
// The event's owner can always cancel it; everyone else needs the cancel-any permission
// (which admins bypass).
func (s *Service) canCancelEvent(ctx context.Context, userId, creatorId uuid.UUID) (bool, error) {
	if userId == creatorId {
		return true, nil
	}
	return s.rbac.IsAuthorized(ctx, userId, rbac.PermissionEventsCancelAny)
}

func (s *Service) CancelEvent(ctx context.Context, id string) error {
	uuidId, err := uuid.Parse(id)
	if err != nil {
		return apperrors.NewBadRequest("invalid event id", err)
	}

	event, err := s.store.getEventById(ctx, uuidId)
	if err != nil {
		return err
	}

	if event.CancelledAt != nil {
		return apperrors.NewBadRequest("event is already cancelled", nil)
	}

	userId := ctx.Value(s.cfg.JWT.User.Key).(*users.Claims).ID

	authorized, err := s.canCancelEvent(ctx, userId, event.CreatedBy)
	if err != nil {
		return err
	}
	if !authorized {
		return apperrors.NewForbidden("you are not allowed to cancel this event")
	}

	return s.store.cancelEvent(ctx, uuidId)
}

func (s *Service) UncancelEvent(ctx context.Context, id string) error {
	uuidId, err := uuid.Parse(id)
	if err != nil {
		return apperrors.NewBadRequest("invalid event id", err)
	}

	event, err := s.store.getEventById(ctx, uuidId)
	if err != nil {
		return err
	}

	if event.CancelledAt == nil {
		return apperrors.NewBadRequest("event is not cancelled", nil)
	}

	userId := ctx.Value(s.cfg.JWT.User.Key).(*users.Claims).ID

	authorized, err := s.canCancelEvent(ctx, userId, event.CreatedBy)
	if err != nil {
		return err
	}
	if !authorized {
		return apperrors.NewForbidden("you are not allowed to restore this event")
	}

	return s.store.uncancelEvent(ctx, uuidId)
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

	claims := ctx.Value(s.cfg.JWT.User.Key).(*users.Claims)

	// RSVP for a kid: it must be the caller's own kid and approved. Ownership is
	// enforced in the query (a kid belonging to someone else is reported as not
	// found), and approval is gated here — never trust the client.
	if req.KidID != nil {
		status, err := s.store.getOwnedKidStatus(ctx, *req.KidID, claims.ID)
		if err != nil {
			return err
		}
		if status != "approved" {
			return apperrors.NewForbidden("this kid is not approved yet")
		}
		return s.store.upsertKidEventResponse(ctx, uuidId, *req.KidID, *req.Response)
	}

	// Own response: only participating accounts may RSVP for themselves. A guardian
	// ("only my kids") account has no own participation.
	if !claims.SelfParticipates {
		return apperrors.NewForbidden("this account does not participate in events itself")
	}
	return s.store.upsertEventResponse(ctx, uuidId, claims.ID, *req.Response)
}

func (s *Service) CreateEvent(ctx context.Context, e Event) (*Event, error) {
	if e.StartTime.Before(time.Now()) {
		return nil, apperrors.NewBadRequest("start time must be in the future", nil)
	}

	claims := ctx.Value(s.cfg.JWT.User.Key).(*users.Claims)
	userId := claims.ID

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

	// Auto-RSVP the creator as going — but only if the account participates itself.
	// (A guardian account can't create events anyway, since it holds no member role.)
	if claims.SelfParticipates {
		if err := s.store.upsertEventResponse(ctx, id, userId, true); err != nil {
			return nil, err
		}
	}

	return &e, nil
}

// ListComments returns the comments on an event, oldest first. Any authenticated
// member may read them; the route middleware enforces authentication.
func (s *Service) ListComments(ctx context.Context, eventId string) ([]*Comment, error) {
	uuidId, err := uuid.Parse(eventId)
	if err != nil {
		return nil, apperrors.NewBadRequest("invalid event id", err)
	}

	if _, err := s.store.getEventById(ctx, uuidId); err != nil {
		return nil, err
	}

	return s.store.getEventComments(ctx, uuidId)
}

// CreateComment adds a comment to an event on behalf of the current account holder.
// Any authenticated member may comment; the body is validated (non-blank, max 500) at
// the API boundary.
func (s *Service) CreateComment(ctx context.Context, eventId string, req CreateCommentRequest) (*Comment, error) {
	uuidId, err := uuid.Parse(eventId)
	if err != nil {
		return nil, apperrors.NewBadRequest("invalid event id", err)
	}

	if _, err := s.store.getEventById(ctx, uuidId); err != nil {
		return nil, err
	}

	userId := ctx.Value(s.cfg.JWT.User.Key).(*users.Claims).ID
	return s.store.createComment(ctx, uuidId, userId, req.Body)
}

// displayName returns the best available display name for the account holder,
// preferring the nick name and falling back to the given name.
func displayName(c *users.Claims) string {
	if c.NickName != nil && *c.NickName != "" {
		return *c.NickName
	}
	if c.GivenName != nil {
		return *c.GivenName
	}
	return ""
}
