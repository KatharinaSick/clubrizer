package events

import (
	"context"

	"github.com/google/uuid"
)

// listener is told about changes to events once they have happened, e.g. to send
// notifications. It only gets the facts — deciding who to tell what is up to the listener.
// Implementations are best-effort and never fail the action that triggered them.
type listener interface {
	EventCancelled(ctx context.Context, change EventChanged)
	EventUncancelled(ctx context.Context, change EventChanged)
	EventDeleted(ctx context.Context, change EventChanged)
	CommentCreated(ctx context.Context, created CommentCreated)
}

// EventChanged describes something done to an event by ActorID.
type EventChanged struct {
	Event   *Event
	ActorID uuid.UUID
	// Attendees are the account holders who said yes, for themselves or one of their kids.
	// Loaded before the change, since deleting an event also deletes its responses.
	Attendees []uuid.UUID
}

// CommentCreated describes a new comment on an event.
type CommentCreated struct {
	Event      *Event
	Comment    *Comment
	AuthorID   uuid.UUID
	AuthorName string
	// Attendees are the account holders who said yes, for themselves or one of their kids.
	Attendees []uuid.UUID
}

// loadAttendees returns who said yes to the event, for a listener. A failure is only logged
// and yields nobody — it must never fail the action itself.
func (s *Service) loadAttendees(ctx context.Context, eventId uuid.UUID) []uuid.UUID {
	attendees, err := s.store.getAttendingUserIds(ctx, eventId)
	if err != nil {
		s.log.Warn("failed to load event attendees", "event", eventId, "error", err)
		return nil
	}
	return attendees
}
