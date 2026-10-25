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
	EventUpdated(ctx context.Context, updated EventUpdated)
	CommentCreated(ctx context.Context, created CommentCreated)
	EventCreated(ctx context.Context, created EventCreated)
	ResponseChanged(ctx context.Context, changed ResponseChanged)
}

// EventChanged describes something done to an event by ActorID.
type EventChanged struct {
	Event   *Event
	ActorID uuid.UUID
	// Attendees are the account holders who said yes, for themselves or one of their kids.
	Attendees []uuid.UUID
}

// EventUpdated describes an edit of an event by ActorID. Before and After are the event as
// it was and as it is now, so the listener can tell what changed.
type EventUpdated struct {
	Before    *Event
	After     *Event
	ActorID   uuid.UUID
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

// EventCreated describes a new event, created by ActorID.
type EventCreated struct {
	Event   *Event
	ActorID uuid.UUID
}

// ResponseChanged describes someone saying yes or no to an event. ActorID is the account
// holder who answered — for themselves or, if KidID is set, for one of their kids. Name is
// the display name of whoever the answer is for.
type ResponseChanged struct {
	Event    *Event
	ActorID  uuid.UUID
	KidID    *uuid.UUID
	Name     string
	Response bool
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
