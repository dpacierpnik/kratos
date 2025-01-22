package detector

import (
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"time"
)

// Detector Detects visits, where an identity is accessed from two different geographic locations (eg countries) within an unreasonably short timeframe, suggesting that conventional travel between those locations would be impossible.
type Detector[T GeoLocated] struct {
	visitsAggregator        VisitsAggregator[T]
	travelDurationEstimator TravelDurationEstimator[T]
}

func NewDetector[T GeoLocated](
	visitsAggregator VisitsAggregator[T],
	travelDurationEstimator TravelDurationEstimator[T],
) *Detector[T] {
	return &Detector[T]{
		visitsAggregator:        visitsAggregator,
		travelDurationEstimator: travelDurationEstimator,
	}
}

type VisitsAggregator[T GeoLocated] interface {
	GetPreviousVisitsAggregate(identityID uuid.UUID) (PreviousVisitsAggregate[T], error)
	AggregateVisit(identityID uuid.UUID, visit Visit[T], detectionResult DetectionResult) error
}

type PreviousVisitsAggregate[T GeoLocated] interface {
	GetLastLegitVisit() (Visit[T], bool)
	IsLegitVisit(currentVisit Visit[T]) bool
}

type TravelDurationEstimator[T GeoLocated] interface {
	EstimateTravelDuration(from T, to T) (time.Duration, error)
}

// DetectImpossibleTravel is a template method for impossible travel detection algorithm. Replaceable parts are extracted to appropriate strategies.
func (d *Detector[T]) DetectImpossibleTravel(identityID uuid.UUID, currentVisit Visit[T]) DetectionResult {
	log.Debugf("Starting impossible travel detection` for visit: `%v`...", currentVisit)
	detectionResult := d.performDetection(identityID, currentVisit)
	saveErr := d.visitsAggregator.AggregateVisit(identityID, currentVisit, detectionResult)
	if saveErr != nil {
		// generally detection succeed, so we still want to mark visit correctly,
		// however data from current visit won't be available in the future detections, so reporting this as a warning
		log.Warnf("Error aggregating current visit, so it won't be available in the future detections: %v", saveErr)
	}
	log.Debugf("Impossible travel detection finished with result `%v` for visit: `%v`", currentVisit, detectionResult)
	return detectionResult
}

func (d *Detector[T]) performDetection(identityID uuid.UUID, currentVisit Visit[T]) DetectionResult {
	previousVisitsAggregate, getErr := d.visitsAggregator.GetPreviousVisitsAggregate(identityID)
	if getErr != nil {
		log.Errorf("Error while getting previous visits aggregate: %v", getErr)
		return DetectionError
	}
	lastLegitVisit, lastLegitVisitFound := previousVisitsAggregate.GetLastLegitVisit()
	if !lastLegitVisitFound {
		// it's a first visit so there is nothing to compare and no reason to not classify visit as legit
		return LegitFirstVisit
	}
	// provides a mechanism to determine if a user's visit qualifies as legitimate basing on history data, without further evaluation of impossible travel conditions.
	// It allows to mitigate the following false-positive cases:
	// - User uses multiple devices and/or multiple applications from multiple IP addresses:
	//   - VPN connection to corporate network with split tunnels:
	//     -- some apps go directly with home IP address
	//     -- some apps (eg using corporate GitHub) go through the VPN in different country
	//
	// - Zoom app may be used on a smartphone:
	//   - IP address can change in seconds between a home Wi-Fi and ISP of a mobile network provider
	//
	// - Foreign travels:
	//   - smartphone can change network between a Wi-Fi on given continent (eg in the office) and ISP of a mobile network provider from a home country connected via roaming
	if previousVisitsAggregate.IsLegitVisit(currentVisit) {
		return LegitVisit
	}
	travel := Travel[T]{
		CurrentVisit:   currentVisit,
		LastLegitVisit: lastLegitVisit,
	}
	if !lastLegitVisit.GeoLocation.IsGeoLocated() || !currentVisit.GeoLocation.IsGeoLocated() {
		log.Warnf("Unable to verify travel duration because of insufficient Geo-location data: %v", travel)
		return InsufficientData
	}
	unreasonablyShort, estimationErr := d.isUnreasonablyShort(travel)
	if estimationErr != nil {
		log.Errorf("Error while estimating travel duration: %v", estimationErr)
		return DetectionError
	}
	if unreasonablyShort {
		return ImpossibleTravel
	}
	return LegitTravel
}

func (d *Detector[T]) isUnreasonablyShort(travel Travel[T]) (bool, error) {
	estimatedDuration, err := d.travelDurationEstimator.EstimateTravelDuration(travel.LastLegitVisit.GeoLocation, travel.CurrentVisit.GeoLocation)
	if err != nil {
		return false, err
	}
	return travel.Duration() < estimatedDuration, nil
}
