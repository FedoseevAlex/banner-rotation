package mab

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUCB1(t *testing.T) {
	var testBannersNum = 10
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
}
