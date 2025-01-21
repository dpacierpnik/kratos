package detector

import (
	"github.com/gofrs/uuid"
	"time"
)

// GeoLocated is marker interface to enforce types which may be Geo-located (like cities, or countries)
type GeoLocated interface {
	IsGeoLocated() bool
}

// Visit Describes user (identity) activity as Geo-located visit. Contains all data required to perform impossible travel detection.
type Visit[T GeoLocated] struct {
	Time        time.Time
	DeviceId    uuid.UUID
	GeoLocation T
}

// Travel describes travel between 2 Geo-located visits.
type Travel[T GeoLocated] struct {
	LastLegitVisit Visit[T]
	CurrentVisit   Visit[T]
}

func (t Travel[T]) Duration() time.Duration {
	return t.LastLegitVisit.Time.Sub(t.CurrentVisit.Time)
}

type DetectionResult string

const (
	LegitFirstVisit  DetectionResult = "LegitFirstVisit"
	LegitVisit       DetectionResult = "LegitVisit"
	LegitTravel      DetectionResult = "LegitTravel"
	InsufficientData DetectionResult = "InsufficientData"
	DetectionError   DetectionResult = "DetectionError"
	ImpossibleTravel DetectionResult = "ImpossibleTravel"
)
