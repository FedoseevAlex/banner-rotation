package mab

import (
	"math"

	"github.com/google/uuid"
)

type BannerData struct {
	ID      uuid.UUID
	GroupID uuid.UUID
	SlotID  uuid.UUID
	Shows   int
	Clicks  int
}

func UCB1(banners []*BannerData, trials int) *BannerData {
	var (
		maxConfidence float64
		bannerToShow  *BannerData
	)

	for _, banner := range banners {
		meanClicks := float64(banner.Clicks) / float64(banner.Shows+1)
		confidence := meanClicks + math.Sqrt(2*math.Log(float64(trials+1))/float64(banner.Shows+1))
		if confidence >= maxConfidence {
			maxConfidence = confidence
			bannerToShow = banner
		}
	}

	bannerToShow.Shows += 1
	return bannerToShow
}
