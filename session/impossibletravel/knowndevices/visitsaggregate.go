package knowndevices

import (
	"github.com/gofrs/uuid"
	"github.com/ory/kratos/session/impossibletravel/detector"
	"reflect"
	"slices"
)

// FIXME Provide more reliable and secure implementation
// FIXME Consider also usage of AuthenticatorAssuranceLevel, to accept visits which are very hard to break (eg recognize visit as legit for AuthenticatorAssuranceLevel >= `aal2`

type VisitsAggregate[T detector.GeoLocated] struct {
	knownDeviceIds []uuid.UUID
	lastLegitVisit detector.Visit[T]
}

func NewVisitsAggregate[T detector.GeoLocated](
	knownDeviceIds []uuid.UUID,
	lastLegitVisit detector.Visit[T],
) *VisitsAggregate[T] {
	return &VisitsAggregate[T]{
		knownDeviceIds: knownDeviceIds,
		lastLegitVisit: lastLegitVisit,
	}
}

func (a *VisitsAggregate[T]) GetLastLegitVisit() (detector.Visit[T], bool) {
	var emptyVisit detector.Visit[T]
	isZero := reflect.DeepEqual(a.lastLegitVisit, emptyVisit)
	return a.lastLegitVisit, !isZero
}

func (a *VisitsAggregate[T]) IsLegitVisit(currentVisit detector.Visit[T]) bool {
	return slices.Contains(a.knownDeviceIds, currentVisit.DeviceId)
}
