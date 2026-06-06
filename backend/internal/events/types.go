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
	Title       string          `db:"title" json:"title" validate:"required"`
	Description string          `db:"description" json:"description"`
	Location    string          `db:"location" json:"location" validate:"required"`
	StartTime   time.Time       `db:"start_time" json:"startTime" validate:"required"`
	CreatedBy   uuid.UUID       `db:"created_by" json:"-"`
	CreatedAt   time.Time       `db:"created_at" json:"-"`
	CategoryID  uuid.UUID       `db:"category" json:"categoryId" validate:"required"`
	Category    Category        `db:"-" json:"category"`
	Creator     Creator         `db:"-" json:"creator"`
	Responses   *EventResponses `db:"-" json:"responses,omitempty"`
}

type EventAttendee struct {
	ID         uuid.UUID `json:"id"`
	GivenName  string    `json:"givenName"`
	FamilyName string    `json:"familyName"`
	NickName   string    `json:"nickName"`
	Picture    *string   `json:"picture"`
	Response   bool      `json:"response"`
}

type EventResponses struct {
	Going               int              `json:"going"`
	NotGoing            int              `json:"notGoing"`
	CurrentUserResponse *bool            `json:"currentUserResponse"`
	Attendees           []*EventAttendee `json:"attendees"`
}

type UpsertEventResponseRequest struct {
	Response *bool `json:"response" validate:"required"`
}
