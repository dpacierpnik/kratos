package crosscountry

import (
	"fmt"
	"time"
)

// FIXME for sake of this example, those are only samples from: https://distancecalculator.globefeed.com/Distance_Between_Countries.asp
var sampleDistances = map[crossCountryTravel]Kilometers{
	{"US", "CA"}: 2_261.43,
	{"JP", "FR"}: 9_849.63,
	{"ZA", "UK"}: 9_880.13,
}

type SpeedBasedTravelsDurationEstimator struct {
	distances map[crossCountryTravel]Kilometers
	speed     KilometersPerHour
}

// NewAvgAirplaneTravelsEstimator creates a SpeedBasedTravelsDurationEstimator with an average flight speed of 567 mph.
func NewAvgAirplaneTravelsEstimator() *SpeedBasedTravelsDurationEstimator {
	// avg airplane speed based on https://distancecalculator.globefeed.com/Distance_Between_Countries.asp
	return NewSampledSpeedBasedTravelsEstimator(MilesPerHour(567).toKilometersPerHour())
}

// NewSpeedOfLightTravelsEstimator creates a SpeedBasedTravelsDurationEstimator with a speed approximation equivalent to the speed of light.
func NewSpeedOfLightTravelsEstimator() *SpeedBasedTravelsDurationEstimator {
	// Approximate value of speed of light in Km / h
	return NewSampledSpeedBasedTravelsEstimator(KilometersPerHour(1_080_000_000))
}

// NewSampledSpeedBasedTravelsEstimator creates a SpeedBasedTravelsDurationEstimator with predefined sample distances and a given speed.
func NewSampledSpeedBasedTravelsEstimator(speed KilometersPerHour) *SpeedBasedTravelsDurationEstimator {
	return NewSpeedBasedTravelsEstimator(sampleDistances, speed)
}

// NewSpeedBasedTravelsEstimator creates a new SpeedBasedTravelsDurationEstimator with provided distances and speed for travel duration calculations.
func NewSpeedBasedTravelsEstimator(distances map[crossCountryTravel]Kilometers, speed KilometersPerHour) *SpeedBasedTravelsDurationEstimator {
	return &SpeedBasedTravelsDurationEstimator{
		distances: distances,
		speed:     speed,
	}
}

func (e *SpeedBasedTravelsDurationEstimator) EstimateTravelDuration(from Country, to Country) (time.Duration, error) {
	distance, err := e.distanceBetween(from, to)
	if err != nil {
		return 0, err
	}
	return e.speed.durationFor(distance), nil
}

func (e *SpeedBasedTravelsDurationEstimator) distanceBetween(from Country, to Country) (Kilometers, error) {
	travel := crossCountryTravel{from, to}
	if distanceInKm, found := e.distances[travel]; found {
		return distanceInKm, nil
	}
	reversedTravel := travel.reverse()
	if distanceInKm, found := e.distances[reversedTravel]; found {
		return distanceInKm, nil
	}
	return 0, fmt.Errorf("no data to estimate distance between %s and %s", from, to)
}

type crossCountryTravel struct {
	from Country
	to   Country
}

func (p crossCountryTravel) reverse() crossCountryTravel {
	return crossCountryTravel{
		from: p.to,
		to:   p.from,
	}
}
