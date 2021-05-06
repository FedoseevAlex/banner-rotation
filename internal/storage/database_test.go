package storage_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/FedoseevAlex/banner-rotation/internal/storage"
	"github.com/FedoseevAlex/banner-rotation/internal/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const dbConnEnvVar = "DB_CONN_STR"

var (
	connectionString = os.Getenv(dbConnEnvVar)
	store            = storage.New(connectionString)
)

type testRotationInfo struct {
	banner types.Banner
	slot   types.Slot
	group  types.Group
}

func TestBanners(t *testing.T) {
	if connectionString == "" {
		t.Skipf("Skipping TestBanners as env var '%s' is not set", dbConnEnvVar)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := store.Connect(ctx)
	require.NoError(t, err)

	banner := types.Banner{ID: uuid.New(), Description: "Some banner"}

	t.Run("check create banner", func(t *testing.T) {
		err := store.AddBanner(ctx, banner)
		require.NoError(t, err)

		dbBanner, err := store.GetBanner(ctx, banner.ID)
		require.NoError(t, err)
		require.Equal(t, banner, dbBanner)
	})

	t.Run("check delete banner", func(t *testing.T) {
		err := store.DeleteBanner(ctx, banner.ID)
		require.NoError(t, err)

		_, err = store.GetBanner(ctx, banner.ID)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestSlots(t *testing.T) {
	if connectionString == "" {
		t.Skipf("Skipping TestSlots as env var '%s' is not set", dbConnEnvVar)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := store.Connect(ctx)
	require.NoError(t, err)

	slot := types.Slot{ID: uuid.New(), Description: "Main slot"}

	t.Run("check create slot", func(t *testing.T) {
		err = store.AddSlot(ctx, slot)
		require.NoError(t, err)

		dbSlot, err := store.GetSlot(ctx, slot.ID)
		require.NoError(t, err)
		require.Equal(t, slot, dbSlot)
	})

	t.Run("check delete slot", func(t *testing.T) {
		err = store.DeleteSlot(ctx, slot.ID)
		require.NoError(t, err)

		_, err = store.GetSlot(ctx, slot.ID)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestGroups(t *testing.T) {
	if connectionString == "" {
		t.Skipf("Skipping TestGroups as env var '%s' is not set", dbConnEnvVar)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := store.Connect(ctx)
	require.NoError(t, err)

	group := types.Group{ID: uuid.New(), Description: "Teenagers"}

	t.Run("check create group", func(t *testing.T) {
		err = store.AddGroup(ctx, group)
		require.NoError(t, err)

		dbGroup, err := store.GetGroup(ctx, group.ID)
		require.NoError(t, err)
		require.Equal(t, group, dbGroup)
	})

	t.Run("check delete group", func(t *testing.T) {
		err = store.DeleteGroup(ctx, group.ID)
		require.NoError(t, err)

		_, err = store.GetGroup(ctx, group.ID)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func createTestRotation(ctx context.Context, t *testing.T, r testRotationInfo) {
	err := store.AddBanner(ctx, r.banner)
	require.NoError(t, err)

	err = store.AddSlot(ctx, r.slot)
	require.NoError(t, err)

	err = store.AddGroup(ctx, r.group)
	require.NoError(t, err)

	err = store.AddRotation(ctx, r.banner.ID, r.slot.ID, r.group.ID)
	require.NoError(t, err)
}

func TestRotations(t *testing.T) {
	if connectionString == "" {
		t.Skipf("Skipping database tests as env var '%s' is not set", dbConnEnvVar)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := store.Connect(ctx)
	require.NoError(t, err)

	rotations := []testRotationInfo{
		{
			banner: types.Banner{ID: uuid.New(), Description: "Some banner"},
			slot:   types.Slot{ID: uuid.New(), Description: "Main slot"},
			group:  types.Group{ID: uuid.New(), Description: "Teenagers"},
		},
		{
			banner: types.Banner{ID: uuid.New(), Description: "Some banner"},
			slot:   types.Slot{ID: uuid.New(), Description: "Main slot"},
			group:  types.Group{ID: uuid.New(), Description: "Teenagers"},
		},
	}

	t.Run("check create rotation", func(t *testing.T) {
		for _, r := range rotations {
			createTestRotation(ctx, t, r)

			dbRotation, err := store.GetRotation(ctx, r.banner.ID, r.slot.ID, r.group.ID)
			require.NoError(t, err)
			require.Equal(t, types.Rotation{BannerID: r.banner.ID, SlotID: r.slot.ID, GroupID: r.group.ID}, dbRotation)
		}
	})

	t.Run("check total shows", func(t *testing.T) {
		testShows := 10
		for _, r := range rotations {
			for i := 0; i < testShows; i++ {
				err := store.AddShow(ctx, r.banner.ID, r.slot.ID, r.group.ID)
				require.NoError(t, err)
			}
		}

		shows, err := store.GetTotalShows(ctx)
		require.NoError(t, err)
		require.Equal(t, int64(len(rotations)*testShows), shows)
	})

	t.Run("check delete rotation1", func(t *testing.T) {
		r := rotations[0]
		err := store.DeleteRotation(ctx, r.banner.ID, r.slot.ID, r.group.ID)
		require.NoError(t, err)
	})

	t.Run("check total shows with deleted rotations", func(t *testing.T) {
		shows, err := store.GetTotalShows(ctx)
		require.NoError(t, err)
		require.Equal(t, int64(10), shows)
	})

	t.Run("check total shows with deleted rotations", func(t *testing.T) {
		shows, err := store.GetTotalShows(ctx)
		require.NoError(t, err)
		require.Equal(t, int64(10), shows)
	})

	t.Run("check rotation mark deleted if banner deleted", func(t *testing.T) {
		testRotation := testRotationInfo{
			banner: types.Banner{ID: uuid.New(), Description: "Some banner"},
			slot:   types.Slot{ID: uuid.New(), Description: "Main slot"},
			group:  types.Group{ID: uuid.New(), Description: "Teenagers"},
		}
		createTestRotation(ctx, t, testRotation)

		err := store.DeleteBanner(ctx, testRotation.banner.ID)
		require.NoError(t, err)

		_, err = store.GetRotation(
			ctx,
			testRotation.banner.ID,
			testRotation.slot.ID,
			testRotation.group.ID,
		)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("check rotation mark deleted if slot deleted", func(t *testing.T) {
		testRotation := testRotationInfo{
			banner: types.Banner{ID: uuid.New(), Description: "Some banner"},
			slot:   types.Slot{ID: uuid.New(), Description: "Main slot"},
			group:  types.Group{ID: uuid.New(), Description: "Teenagers"},
		}
		createTestRotation(ctx, t, testRotation)

		err := store.DeleteSlot(ctx, testRotation.slot.ID)
		require.NoError(t, err)

		_, err = store.GetRotation(
			ctx,
			testRotation.banner.ID,
			testRotation.slot.ID,
			testRotation.group.ID,
		)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("check rotation mark deleted if group deleted", func(t *testing.T) {
		testRotation := testRotationInfo{
			banner: types.Banner{ID: uuid.New(), Description: "Some banner"},
			slot:   types.Slot{ID: uuid.New(), Description: "Main slot"},
			group:  types.Group{ID: uuid.New(), Description: "Teenagers"},
		}
		createTestRotation(ctx, t, testRotation)

		err := store.DeleteGroup(ctx, testRotation.group.ID)
		require.NoError(t, err)

		_, err = store.GetRotation(
			ctx,
			testRotation.banner.ID,
			testRotation.slot.ID,
			testRotation.group.ID,
		)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}
