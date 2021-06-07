package app

import (
	"context"

	"github.com/FedoseevAlex/banner-rotation/internal/config"
	"github.com/FedoseevAlex/banner-rotation/internal/logger"
	"github.com/FedoseevAlex/banner-rotation/internal/rotators/mab"
	"github.com/FedoseevAlex/banner-rotation/internal/storage"
	"github.com/FedoseevAlex/banner-rotation/internal/types"
	"github.com/google/uuid"
)

func New(config config.Config) (*App, error) {
	log := logger.New(config.Log.Level, config.Log.File)

	store := storage.New(config.Storage.DBConnectionString)
	err := store.Connect()
	if err != nil {
		log.Error(
			"failed to connect to storage",
			types.LogFields{
				"error": err,
			},
		)
		return nil, err
	}

	rotator := &mab.MultiArmedBandit{}
	log.Debug(
		"Application created successfully",
		types.LogFields{},
	)
	return &App{Rotator: rotator, Storage: store, Log: log}, nil
}

type App struct {
	Rotator types.Rotator
	Storage types.Storager
	Log     types.Logger
}

func (a *App) AddBanner(ctx context.Context, description string) (types.Banner, error) {
	bannerID, err := uuid.NewRandom()
	if err != nil {
		a.Log.Error(
			"failed to create random uuid for banner",
			types.LogFields{"error": err},
		)
		return types.Banner{}, err
	}

	banner := types.Banner{
		ID:          bannerID,
		Description: description,
	}

	err = a.Storage.AddBanner(ctx, banner)
	if err != nil {
		a.Log.Error(
			"failed to add banner",
			types.LogFields{"error": err},
		)
		return types.Banner{}, err
	}

	a.Log.Trace(
		"add banner",
		types.LogFields{
			"banner_id": bannerID.String(),
		},
	)

	return banner, nil
}

func (a *App) DeleteBanner(ctx context.Context, bannerID uuid.UUID) error {
	err := a.Storage.DeleteBanner(ctx, bannerID)
	if err != nil {
		a.Log.Error(
			"failed to delete banner",
			types.LogFields{"error": err},
		)
		return err
	}
	return nil
}

func (a *App) GetBanner(ctx context.Context, bannerID uuid.UUID) (types.Banner, error) {
	banner, err := a.Storage.GetBanner(ctx, bannerID)
	if err != nil {
		a.Log.Error(
			"failed to get banner from database",
			types.LogFields{"error": err},
		)
		return types.Banner{}, err
	}
	return banner, nil
}

func (a *App) AddSlot(ctx context.Context, description string) (types.Slot, error) {
	slotID, err := uuid.NewRandom()
	if err != nil {
		a.Log.Error(
			"failed to create random uuid for slot",
			types.LogFields{"error": err},
		)
		return types.Slot{}, err
	}

	slot := types.Slot{
		ID:          slotID,
		Description: description,
	}

	err = a.Storage.AddSlot(ctx, slot)
	if err != nil {
		a.Log.Error(
			"failed to add slot",
			types.LogFields{"error": err},
		)
		return types.Slot{}, err
	}

	a.Log.Trace(
		"add slot",
		types.LogFields{
			"banner_id": slotID.String(),
		},
	)

	return slot, nil
}

func (a *App) DeleteSlot(ctx context.Context, slotID uuid.UUID) error {
	err := a.Storage.DeleteSlot(ctx, slotID)
	if err != nil {
		a.Log.Error(
			"failed to delete slot",
			types.LogFields{"error": err},
		)
		return err
	}
	return nil
}

func (a *App) GetSlot(ctx context.Context, slotID uuid.UUID) (types.Slot, error) {
	slot, err := a.Storage.GetSlot(ctx, slotID)
	if err != nil {
		a.Log.Error(
			"failed to get slot from database",
			types.LogFields{"error": err},
		)
		return types.Slot{}, err
	}
	return slot, nil
}

func (a *App) AddGroup(ctx context.Context, description string) (types.Group, error) {
	groupID, err := uuid.NewRandom()
	if err != nil {
		a.Log.Error(
			"failed to create random uuid for group",
			types.LogFields{"error": err},
		)
		return types.Group{}, err
	}

	group := types.Group{
		ID:          groupID,
		Description: description,
	}

	err = a.Storage.AddGroup(ctx, group)
	if err != nil {
		a.Log.Error(
			"failed to add group",
			types.LogFields{"error": err},
		)
		return types.Group{}, err
	}

	a.Log.Trace(
		"add group",
		types.LogFields{
			"banner_id": groupID.String(),
		},
	)

	return group, nil
}

func (a *App) DeleteGroup(ctx context.Context, groupID uuid.UUID) error {
	err := a.Storage.DeleteGroup(ctx, groupID)
	if err != nil {
		a.Log.Error(
			"failed to delete group",
			types.LogFields{"error": err},
		)
		return err
	}
	return nil
}

func (a *App) GetGroup(ctx context.Context, groupID uuid.UUID) (types.Group, error) {
	group, err := a.Storage.GetGroup(ctx, groupID)
	if err != nil {
		a.Log.Error(
			"failed to get group from database",
			types.LogFields{"error": err},
		)
		return types.Group{}, err
	}
	return group, nil
}

func (a *App) AddRotation(ctx context.Context, bannerID, slotID, groupID uuid.UUID) (types.Rotation, error) {
	rotation, err := a.Storage.AddRotation(ctx, bannerID, slotID, groupID)
	if err != nil {
		a.Log.Error(
			"failed to add new rotation",
			types.LogFields{
				"error":     err,
				"banner_id": bannerID.String(),
				"slot_id":   slotID.String(),
				"group_id":  groupID.String(),
			},
		)
		return types.Rotation{}, err
	}
	return rotation, nil
}

func (a *App) GetRotation(ctx context.Context, bannerID, slotID, groupID uuid.UUID) (types.Rotation, error) {
	rotation, err := a.Storage.GetRotation(ctx, bannerID, slotID, groupID)
	if err != nil {
		a.Log.Error(
			"failed to get rotation",
			types.LogFields{
				"error":     err,
				"banner_id": bannerID.String(),
				"slot_id":   slotID.String(),
				"group_id":  groupID.String(),
			},
		)
		return types.Rotation{}, err
	}
	return rotation, nil
}

func (a *App) DeleteRotation(ctx context.Context, bannerID, slotID, groupID uuid.UUID) error {
	err := a.Storage.DeleteRotation(ctx, bannerID, slotID, groupID)
	if err != nil {
		a.Log.Error(
			"failed to delete rotation",
			types.LogFields{
				"error":     err,
				"banner_id": bannerID.String(),
				"slot_id":   slotID.String(),
				"group_id":  groupID.String(),
			},
		)
		return err
	}
	return nil
}

func (a *App) RegisterClick(ctx context.Context, bannerID, slotID, groupID uuid.UUID) error {
	err := a.Storage.AddClick(ctx, bannerID, slotID, groupID)
	if err != nil {
		a.Log.Error(
			"failed to register click for rotation",
			types.LogFields{
				"error":     err,
				"banner_id": bannerID.String(),
				"slot_id":   slotID.String(),
				"group_id":  groupID.String(),
			},
		)
		return err
	}
	return nil
}

func (a *App) ChooseBanner(ctx context.Context, slotID, groupID uuid.UUID) (types.Rotation, error) {
	a.Log.Debug(
		"choose banner",
		types.LogFields{
			"slot_id":  slotID.String(),
			"group_id": groupID.String(),
		},
	)

	trials, err := a.Storage.GetTotalShows(ctx)
	if err != nil {
		a.Log.Error(
			"failed to fetch total trials from db",
			types.LogFields{
				"error": err,
			},
		)
		return types.Rotation{}, err
	}

	rotations, err := a.Storage.GetAllRotations(ctx)
	if err != nil {
		a.Log.Error(
			"failed to fetch all rotations from db",
			types.LogFields{
				"error": err,
			},
		)
		return types.Rotation{}, err
	}

	a.Rotator.Load(rotations, trials)
	rotationToShow := a.Rotator.Rotate()

	// Register show for rotation
	err = a.Storage.AddShow(
		ctx,
		rotationToShow.BannerID,
		rotationToShow.SlotID,
		rotationToShow.GroupID,
	)
	if err != nil {
		a.Log.Error(
			"failed to register show for rotation",
			types.LogFields{
				"error":     err,
				"banner_id": rotationToShow.BannerID,
				"slot_id":   rotationToShow.SlotID,
				"group_id":  rotationToShow.GroupID,
			},
		)
		return types.Rotation{}, err
	}

	a.Log.Debug(
		"rotation to show has been chosen",
		types.LogFields{
			"banner_id": rotationToShow.BannerID,
			"slot_id":   rotationToShow.SlotID,
			"group_id":  rotationToShow.GroupID,
			"shows":     rotationToShow.Shows,
			"clicks":    rotationToShow.Clicks,
		},
	)

	return rotationToShow, nil
}

func (a *App) GetStats(ctx context.Context, bannerID, slotID, groupID uuid.UUID) ([]types.Event, error) {
	a.Log.Debug(
		"get statistics",
		types.LogFields{
			"banner_id": bannerID.String(),
			"slot_id":   slotID.String(),
			"group_id":  groupID.String(),
		},
	)

	events, err := a.Storage.GetRotationStats(ctx, bannerID, slotID, groupID)
	if err != nil {
		a.Log.Error(
			"failed to get rotation stats",
			types.LogFields{
				"error":     err,
				"banner_id": bannerID.String(),
				"slot_id":   slotID.String(),
				"group_id":  groupID.String(),
			},
		)
		return nil, err
	}

	return events, nil
}

func (a *App) GetLogger(name string) types.Logger {
	return a.Log.ChildLogger(name)
}
