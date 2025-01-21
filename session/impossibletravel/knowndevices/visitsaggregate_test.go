package knowndevices

import (
	"github.com/gofrs/uuid"
	"github.com/ory/kratos/session/impossibletravel/detector"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewVisitsAggregate(t *testing.T) {

	t.Run("should return last legit visit", func(t *testing.T) {

		// given
		lastLegitVisit := detector.Visit[anythingGeoLocated]{
			Time:        time.Now(),
			DeviceId:    uuid.Must(uuid.NewV4()),
			GeoLocation: "Location",
		}
		visitsAggregateUT := NewVisitsAggregate(make([]uuid.UUID, 0), lastLegitVisit)

		// when
		actualLastLegitVisit, ok := visitsAggregateUT.GetLastLegitVisit()

		// then
		require.True(t, ok)
		assert.Equal(t, lastLegitVisit, actualLastLegitVisit)
	})

	t.Run("should return empty last legit visit", func(t *testing.T) {

		// given
		lastLegitVisit := detector.Visit[anythingGeoLocated]{}
		visitsAggregateUT := NewVisitsAggregate(make([]uuid.UUID, 0), lastLegitVisit)

		// when
		actualLastLegitVisit, ok := visitsAggregateUT.GetLastLegitVisit()

		// then
		require.False(t, ok)
		assert.Equal(t, lastLegitVisit, actualLastLegitVisit)
	})

	t.Run("should recognize visit as legit if device is known", func(t *testing.T) {

		// given
		lastLegitVisit := detector.Visit[anythingGeoLocated]{}
		currentVisit := detector.Visit[anythingGeoLocated]{
			Time:        time.Now(),
			DeviceId:    uuid.Must(uuid.NewV4()),
			GeoLocation: "Location",
		}
		knownDevices := []uuid.UUID{currentVisit.DeviceId}
		visitsAggregateUT := NewVisitsAggregate(knownDevices, lastLegitVisit)

		// when
		result := visitsAggregateUT.IsLegitVisit(currentVisit)

		// then
		require.True(t, result)
	})

	t.Run("should not recognize visit as legit if device is unknown", func(t *testing.T) {

		// given
		lastLegitVisit := detector.Visit[anythingGeoLocated]{}
		currentVisit := detector.Visit[anythingGeoLocated]{
			Time:        time.Now(),
			DeviceId:    uuid.Must(uuid.NewV4()),
			GeoLocation: "Location",
		}
		var knownDevices []uuid.UUID
		visitsAggregateUT := NewVisitsAggregate(knownDevices, lastLegitVisit)

		// when
		result := visitsAggregateUT.IsLegitVisit(currentVisit)

		// then
		require.False(t, result)
	})
}

type anythingGeoLocated string

func (l anythingGeoLocated) IsGeoLocated() bool {
	return l != ""
}
