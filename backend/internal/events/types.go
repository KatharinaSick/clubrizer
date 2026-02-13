package events

import (
	"github.com/google/uuid"
	"time"
)

type Category struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Color     string    `db:"color" json:"color"`
	Picture   string    `db:"picture" json:"picture"`
	CreatedBy string    `db:"created_by" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"-"`
}

type Event struct {
	ID          uuid.UUID `db:"id" json:"-"`
	Title       string    `db:"title" json:"title" validate:"required"`
	Description string    `db:"description" json:"description"`
	Location    string    `db:"location" json:"location" validate:"required"`
	StartTime   time.Time `db:"start_time" json:"startTime" validate:"required"`
	CreatedBy   uuid.UUID `db:"created_by" json:"-"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
	CategoryID  uuid.UUID `db:"category" json:"categoryId" validate:"required"`
	Category    Category  `db:"-" json:"category"`
}
