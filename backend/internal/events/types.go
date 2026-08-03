package events

import (
	"github.com/google/uuid"
	"time"
)

type Category struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Color       string    `db:"color" json:"color"`
	Picture     string    `db:"picture" json:"picture"`
	SortOrder   int       `db:"sort_order" json:"sortOrder"`
	CustomLabel *string   `db:"custom_label" json:"customLabel,omitempty"`
	CanCreate   *bool     `db:"-" json:"canCreate,omitempty"`
	CreatedBy   string    `db:"created_by" json:"-"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
}

type Creator struct {
	ID         uuid.UUID `json:"id"`
	GivenName  string    `json:"givenName"`
	FamilyName string    `json:"familyName"`
	NickName   string    `json:"nickName"`
	Picture    *string   `json:"picture"`
}

type Event struct {
	ID          uuid.UUID       `db:"id" json:"id"`
	Title       string          `db:"title" json:"title" validate:"required,notblank"`
	Description string          `db:"description" json:"description"`
	Location    string          `db:"location" json:"location" validate:"required,notblank"`
	StartTime   time.Time       `db:"start_time" json:"startTime" validate:"required"`
	CreatedBy   uuid.UUID       `db:"created_by" json:"-"`
	CreatedAt   time.Time       `db:"created_at" json:"-"`
	CategoryID  uuid.UUID       `db:"category" json:"categoryId" validate:"required"`
	Category    Category        `db:"-" json:"category"`
	Creator     Creator         `db:"-" json:"creator"`
	CancelledAt *time.Time      `db:"cancelled_at" json:"cancelledAt,omitempty"`
	Responses   *EventResponses `db:"-" json:"responses,omitempty"`
	CanDelete   *bool           `db:"-" json:"canDelete,omitempty"`
	CanCancel   *bool           `db:"-" json:"canCancel,omitempty"`
}

type EventAttendee struct {
	ID         uuid.UUID `json:"id"`
	GivenName  string    `json:"givenName"`
	FamilyName string    `json:"familyName"`
	NickName   string    `json:"nickName"`
	Picture    *string   `json:"picture"`
	Response   bool      `json:"response"`
	IsKid      bool      `json:"isKid"`
	// Parent is the owning parent's display name, set only for kid attendees so the
	// UI can show "kid of {parent}".
	Parent *string `json:"parent,omitempty"`
}

// MyResponse is one person the current account can RSVP for on this event (the
// account holder — only if they participate — plus each approved kid), with that
// person's current response (nil = not yet responded). Drives the per-person RSVP UI.
type MyResponse struct {
	ID       uuid.UUID `json:"id"`
	IsSelf   bool      `json:"isSelf"`
	Name     string    `json:"name"`
	Picture  *string   `json:"picture"`
	Response *bool     `json:"response"`
}

type EventResponses struct {
	Going       int              `json:"going"`
	NotGoing    int              `json:"notGoing"`
	Attendees   []*EventAttendee `json:"attendees"`
	MyResponses []*MyResponse    `json:"myResponses"`
}

// responder identifies the current account when computing its per-person RSVP list.
type responder struct {
	UserID           uuid.UUID
	Name             string
	Picture          *string
	SelfParticipates bool
}

// UpsertEventResponseRequest sets a response for one person. KidID nil means the
// account holder's own response; KidID set means that kid's response.
type UpsertEventResponseRequest struct {
	Response *bool      `json:"response" validate:"required"`
	KidID    *uuid.UUID `json:"kidId"`
}

// Comment is a message an account holder left on an event, with the author's display
// details joined in so the UI can render it without a second lookup.
type Comment struct {
	ID        uuid.UUID `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	Author    Creator   `json:"author"`
}

// CreateCommentRequest is the body of a new comment. The 500-char cap mirrors the
// event_comments.body column length; notblank rejects whitespace-only messages.
type CreateCommentRequest struct {
	Body string `json:"body" validate:"required,notblank,max=500"`
}
