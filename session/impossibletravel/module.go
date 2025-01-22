package impossibletravel

import (
	"github.com/ory/kratos/session/impossibletravel/crosscountry"
	"github.com/ory/kratos/session/impossibletravel/detector"
	"github.com/ory/kratos/session/impossibletravel/knowndevices"
)

type Config struct {
	// FIXME define config which allows to create wanted detector
}

func NewDetectorFrom(config Config) *detector.Detector[crosscountry.Country] {
	// FIXME use configuration to create detector and it's deps
	visitsAggregator := knowndevices.NewVisitsAggregator[crosscountry.Country]()
	travelDurationEstimator := crosscountry.NewAvgAirplaneTravelsEstimator()
	return detector.NewDetector(visitsAggregator, travelDurationEstimator)
}
