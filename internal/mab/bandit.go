package mab

import (
	"math"

	"github.com/google/uuid"
)

type BannerData struct {
	ID          uuid.UUID
	GroupID     uuid.UUID
	SlotID      uuid.UUID
	Shows       int
	Clicks      int
	Description string
}

type UserGroup struct {
	ID          uuid.UUID
	Description string
}

type Slot struct {
	ID          uuid.UUID
	Description string
}

func UCB1(banners []*BannerData, trials int) *BannerData {
	var (
		maxConfidence float64
		bannerToShow  *BannerData
	)

	if trials == 0 {
		trials = 1
	}

	for _, banner := range banners {
		shows := banner.Shows
		if shows == 0 {
			shows = 1
		}

		meanClicks := float64(banner.Clicks) / float64(shows)
		confidence := meanClicks + math.Sqrt(2*math.Log(float64(trials))/float64(shows))
		if confidence >= maxConfidence {
			maxConfidence = confidence
			bannerToShow = banner
		}
	}

	bannerToShow.Shows += 1
	return bannerToShow
}
