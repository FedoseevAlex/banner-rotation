package storage

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const dbConnEnvVar = "DB_CONN_STR"

func TestDatabase(t *testing.T) {
	connectionString := os.Getenv(dbConnEnvVar)

	if connectionString == "" {
		t.Skipf("Skipping database tests as env var '%s' is not set", dbConnEnvVar)
	}

	store := New(connectionString)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := store.Connect(ctx)
	require.NoError(t, err)

	banner := Banner{ID: uuid.New(), Description: "Some banner"}
	slot := Slot{ID: uuid.New(), Description: "Main slot"}
	group := Group{ID: uuid.New(), Description: "Teenagers"}

	t.Run("check create banner slot and group", func(t *testing.T) {
		err := store.AddBanner(ctx, banner)
		require.NoError(t, err)

		dbBanner, err := store.GetBanner(ctx, banner.ID)
		require.NoError(t, err)
		require.Equal(t, banner, dbBanner)

		err = store.AddSlot(ctx, slot)
		require.NoError(t, err)

		dbSlot, err := store.GetSlot(ctx, slot.ID)
		require.NoError(t, err)
		require.Equal(t, slot, dbSlot)

		err = store.AddGroup(ctx, group)
		require.NoError(t, err)

		dbGroup, err := store.GetGroup(ctx, group.ID)
		require.NoError(t, err)
		require.Equal(t, group, dbGroup)

	})

	t.Run("check create rotation", func(t *testing.T) {
		err := store.AddRotation(ctx, banner.ID, slot.ID, group.ID)
		require.NoError(t, err)

		dbRotation, err := store.GetRotation(ctx, banner.ID, slot.ID, group.ID)
		require.NoError(t, err)
		require.Equal(t, Rotation{BannerID: banner.ID, SlotID: slot.ID, GroupID: group.ID}, dbRotation)
	})

	t.Run("check total shows", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			err := store.AddShow(ctx, banner.ID, slot.ID, group.ID)
			require.NoError(t, err)
		}

		shows, err := store.GetTotalShows(ctx)
		require.NoError(t, err)
		require.Equal(t, int64(10), shows)
	})

	t.Run("check delete rotation", func(t *testing.T) {
		err := store.DeleteRotation(ctx, banner.ID, slot.ID, group.ID)
		require.NoError(t, err)
	})

	t.Run("check delete cascade", func(t *testing.T) {
		err := store.AddRotation(ctx, banner.ID, slot.ID, group.ID)
		require.NoError(t, err)

		dbRotation, err := store.GetRotation(ctx, banner.ID, slot.ID, group.ID)
		require.NoError(t, err)
		require.Equal(t, Rotation{BannerID: banner.ID, SlotID: slot.ID, GroupID: group.ID}, dbRotation)

		err = store.DeleteBanner(ctx, banner.ID)
		require.NoError(t, err)

		_, err = store.GetRotation(ctx, banner.ID, slot.ID, group.ID)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

}
