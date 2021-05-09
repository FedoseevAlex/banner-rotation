package types

import (
	"context"

	"github.com/google/uuid"
)

type Banner struct {
	ID          uuid.UUID
	Description string
}

type Slot struct {
	ID          uuid.UUID
	Description string
}

type Group struct {
	ID          uuid.UUID
	Description string
}

type Rotation struct {
	BannerID uuid.UUID
	SlotID   uuid.UUID
	GroupID  uuid.UUID
	Shows    int
	Clicks   int
}

type Storager interface {
	Connect(ctx context.Context) (err error)
	Close(ctx context.Context) error
	// Banner operations
	AddBanner(ctx context.Context, banner Banner) error
	GetBanner(ctx context.Context, bannerID uuid.UUID) (Banner, error)
	DeleteBanner(ctx context.Context, bannerID uuid.UUID) error
	// Slot operations
	AddSlot(ctx context.Context, slot Slot) error
	GetSlot(ctx context.Context, slotID uuid.UUID) (Slot, error)
	DeleteSlot(ctx context.Context, slotID uuid.UUID) error
	// Group operations
	AddGroup(ctx context.Context, group Group) error
	GetGroup(ctx context.Context, groupID uuid.UUID) (Group, error)
	DeleteGroup(ctx context.Context, groupID uuid.UUID) error
	// Rotation operations
	AddRotation(ctx context.Context, bannerID, slotID, groupID uuid.UUID) error
	DeleteRotation(ctx context.Context, bannerID, slotID, groupID uuid.UUID) error
	GetRotation(ctx context.Context, bannerID, slotID, groupID uuid.UUID) (Rotation, error)
	AddShow(ctx context.Context, bannerID, slotID, groupID uuid.UUID) error
	AddClick(ctx context.Context, bannerID, slotID, groupID uuid.UUID) error
	GetAllRotations(ctx context.Context) ([]Rotation, error)
	// Get total amount of shows
	GetTotalShows(ctx context.Context) (totalShows int64, err error)
}

type Logger interface {
	Debug(msg string, args ...map[string]interface{})
	Info(msg string, args ...map[string]interface{})
	Warn(msg string, args ...map[string]interface{})
	Error(msg string, args ...map[string]interface{})
	Trace(msg string, args ...map[string]interface{})
}

type Rotator interface {
	Rotate() Rotation
	Load(rotations []Rotation, trials int64)
}

type Application interface {
	AddBanner(ctx context.Context, description string) (Banner, error)
	DeleteBanner(ctx context.Context, bannerID uuid.UUID) error
	GetBanner(ctx context.Context, bannerID uuid.UUID) (Banner, error)

	AddSlot(ctx context.Context, description string) (Slot, error)
	DeleteSlot(ctx context.Context, slotID uuid.UUID) error
	GetSlot(ctx context.Context, slotID uuid.UUID) (Slot, error)

	AddGroup(ctx context.Context, description string) (Group, error)
	DeleteGroup(ctx context.Context, groupID uuid.UUID) error
	GetGroup(ctx context.Context, groupID uuid.UUID) (Group, error)

	RegisterClick(ctx context.Context, bannerID, slotID, groupID uuid.UUID) error
	ChooseBanner(ctx context.Context, slotID, groupID uuid.UUID) (Rotation, error)

	GetLogger(name string) Logger
}
