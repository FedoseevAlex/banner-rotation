package mab

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUCB1(t *testing.T) {
	testBannersNum := 10
	banners := make([]*BannerData, 0, testBannersNum)

	for i := 0; i < testBannersNum; i++ {
		bannerID := uuid.New()
		banner := BannerData{
			ID:      bannerID,
			GroupID: uuid.New(),
			SlotID:  uuid.New(),
		}
		banners = append(banners, &banner)
	}

	t.Run("check all banners are shown", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			UCB1(banners, i)
		}

		// Extract Shows property
		allShows := make([]int, 0, testBannersNum)
		for _, banner := range banners {
			allShows = append(allShows, banner.Shows)
		}

		require.NotContains(t, allShows, 0)
	})

	t.Run("check popular banner have maximum shows", func(t *testing.T) {
		// Add "popular" banner
		popularID := uuid.New()
		banners = append(banners, &BannerData{
			ID:      popularID,
			GroupID: uuid.New(),
			SlotID:  uuid.New(),
			Clicks:  10,
		})

		for i := 0; i < 100; i++ {
			UCB1(banners, i)
		}

		// Find most displayed banner
		var (
			maxShows      int
			popularBanner *BannerData
		)
		for _, banner := range banners {
			if banner.Shows > maxShows {
				popularBanner = banner
			}
		}

		require.Equal(t, popularID, popularBanner.ID)
	})
}
