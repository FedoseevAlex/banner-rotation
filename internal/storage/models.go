package storage

import (
	"database/sql"

	"github.com/google/uuid"
)

type banner struct {
	ID          uuid.UUID    `db:"id"`
	Description string       `db:"description"`
	Deleted     bool         `db:"deleted"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
}

type slot struct {
	ID          uuid.UUID    `db:"id"`
	Description string       `db:"description"`
	Deleted     bool         `db:"deleted"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
}

type group struct {
	ID          uuid.UUID    `db:"id"`
	Description string       `db:"description"`
	Deleted     bool         `db:"deleted"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
}

type rotation struct {
	BannerID  uuid.UUID    `db:"banner_id"`
	SlotID    uuid.UUID    `db:"slot_id"`
	GroupID   uuid.UUID    `db:"group_id"`
	Shows     int          `db:"shows"`
	Clicks    int          `db:"clicks"`
	Deleted   bool         `db:"deleted"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}
