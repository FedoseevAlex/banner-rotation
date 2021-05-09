package mab

import (
	"math"

	"github.com/FedoseevAlex/banner-rotation/internal/types"
)

type MultiArmedBandit struct {
	BannersDatas []types.Rotation
	Trials       int64
}

func (mab *MultiArmedBandit) Rotate() (rotation types.Rotation) {
	return UCB1(mab.BannersDatas, mab.Trials)
}

func (mab *MultiArmedBandit) Load(rotations []types.Rotation, trials int64) {
	mab.Trials = trials
	mab.BannersDatas = rotations
}

func UCB1(rotations []types.Rotation, trials int64) types.Rotation {
	var (
		maxConfidence  float64
		rotationToShow types.Rotation
	)

	for _, rotation := range rotations {
		meanClicks := float64(rotation.Clicks) / float64(rotation.Shows+1)
		confidence := meanClicks + math.Sqrt(2*math.Log(float64(trials+1))/float64(rotation.Shows+1))
		if confidence >= maxConfidence {
			maxConfidence = confidence
			rotationToShow = rotation
		}
	}

	return rotationToShow
}
