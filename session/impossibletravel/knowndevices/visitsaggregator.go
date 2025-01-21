package knowndevices

import (
	"github.com/gofrs/uuid"
	"github.com/ory/kratos/session/impossibletravel/detector"
	"slices"
)

type VisitsAggregator[T detector.GeoLocated] struct {
	fakeStorage map[uuid.UUID]*VisitsAggregate[T]
}

func NewVisitsAggregator[T detector.GeoLocated]() *VisitsAggregator[T] {
	return &VisitsAggregator[T]{}
}

func (a *VisitsAggregator[T]) GetPreviousVisitsAggregate(identityID uuid.UUID) (detector.PreviousVisitsAggregate[T], error) {
	// FIXME read data stored for given Identity
	visitsAggregate := a.fakeStorage[identityID]
	return visitsAggregate, nil
}

func (a *VisitsAggregator[T]) AggregateVisit(identityID uuid.UUID, visit detector.Visit[T], detectionResult detector.DetectionResult) error {
	// FIXME read data stored for given Identity and add new visit information
	//  - preferable atomic/transactional operation so aggregates are not messed up in case of many simultaneous activities
	//  - however maybe it's worth to consider also, that those are not very crucial data and loosing a single visit in aggregate may not be very harmful
	// silly implementation just for sake of an example
	if isLegit(detectionResult) {
		aggregate, found := a.fakeStorage[identityID]
		if found {
			aggregate.lastLegitVisit = visit
			// only unique entries
			if !slices.Contains(aggregate.knownDeviceIds, visit.DeviceId) {
				aggregate.knownDeviceIds = append(aggregate.knownDeviceIds, visit.DeviceId)
			}
		} else {
			aggregate = NewVisitsAggregate(make([]uuid.UUID, 0), visit)
			a.fakeStorage[identityID] = aggregate
		}
	}
	return nil
}

func isLegit(detectionResult detector.DetectionResult) bool {
	switch detectionResult {
	case detector.LegitFirstVisit:
	case detector.LegitVisit:
	case detector.LegitTravel:
		return true
	}
	return false
}
