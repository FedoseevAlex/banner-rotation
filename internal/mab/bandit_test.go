package mab

import (
	"testing"

	"github.com/FedoseevAlex/banner-rotation/internal/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUCB1(t *testing.T) {
	testBannersNum := 10
	rotations := make([]types.Rotation, 0, testBannersNum)

	for i := 0; i < testBannersNum; i++ {
		bannerID := uuid.New()
		banner := types.Rotation{
			BannerID: bannerID,
			GroupID:  uuid.New(),
			SlotID:   uuid.New(),
		}
		rotations = append(rotations, banner)
	}

	t.Run("check all banners are shown", func(t *testing.T) {
		rotationsShows := make(map[types.Rotation]int)

		for i := 0; i < 100; i++ {
			rotation := UCB1(rotations, i)
			rotationsShows[rotation]++
		}

		var shows []int
		for _, show := range rotationsShows {
			shows = append(shows, show)
		}

		require.NotContains(t, shows, 0)
	})

	t.Run("check popular banner have maximum shows", func(t *testing.T) {
		// Add "popular" banner
		popularID := uuid.New()
		rotations = append(rotations, types.Rotation{
			BannerID: popularID,
			GroupID:  uuid.New(),
			SlotID:   uuid.New(),
			Clicks:   10,
		})

		rotationsShows := make(map[types.Rotation]int)
		for i := 0; i < 100; i++ {
			rotation := UCB1(rotations, i)
			rotationsShows[rotation]++
		}

		// Find most displayed banner
		var (
			maxShows      int
			popularBanner types.Rotation
		)

		for rotation, shows := range rotationsShows {
			if shows > maxShows {
				popularBanner = rotation
			}
		}

		require.Equal(t, popularID, popularBanner.BannerID)
	})
}
