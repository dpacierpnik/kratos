package crosscountry

import "time"

type Country string

func (c Country) IsGeoLocated() bool {
	return c != ""
}

type Kilometers float64
type KilometersPerHour float64

func (kph KilometersPerHour) durationFor(distance Kilometers) time.Duration {
	durationInHours := float64(distance) / float64(kph)
	return time.Duration(durationInHours) * time.Hour
}

type MilesPerHour float64

func (mph MilesPerHour) toKilometersPerHour() KilometersPerHour {
	return KilometersPerHour(float64(mph) * 1.609344)
}
