package storage

import "github.com/google/uuid"

type banner struct {
	ID          uuid.UUID `db:"id"`
	Description string    `db:"description"`
}

type slot struct {
	ID          uuid.UUID `db:"id"`
	Description string    `db:"description"`
}

type group struct {
	ID          uuid.UUID `db:"id"`
	Description string    `db:"description"`
}

type rotation struct {
	BannerID uuid.UUID `db:"banner_id"`
	SlotID   uuid.UUID `db:"slot_id"`
	GroupID  uuid.UUID `db:"group_id"`
	Shows    int       `db:"shows"`
	Clicks   int       `db:"clicks"`
}
