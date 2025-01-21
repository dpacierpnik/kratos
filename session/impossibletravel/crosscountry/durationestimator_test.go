package crosscountry

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSpeedBasedTravelsDurationEstimator(t *testing.T) {

	t.Run("should correctly estimate time based on distance and speed", func(t *testing.T) {

		// given
		distances := map[crossCountryTravel]Kilometers{
			crossCountryTravel{from: "A", to: "B"}: 100,
		}
		speed := KilometersPerHour(10)

		estimatorUT := NewSpeedBasedTravelsEstimator(distances, speed)

		// when
		estimation, err := estimatorUT.EstimateTravelDuration("A", "B")

		// then
		require.NoError(t, err)
		assert.Equal(t, 10*time.Hour, estimation)
	})

	t.Run("should correctly estimate time based on distance and speed for reverse direction", func(t *testing.T) {

		// given
		distances := map[crossCountryTravel]Kilometers{
			crossCountryTravel{from: "A", to: "B"}: 100,
		}
		speed := KilometersPerHour(10)

		estimatorUT := NewSpeedBasedTravelsEstimator(distances, speed)

		// when
		estimation, err := estimatorUT.EstimateTravelDuration("B", "A")

		// then
		require.NoError(t, err)
		assert.Equal(t, 10*time.Hour, estimation)
	})

	t.Run("should return error for unknown countries", func(t *testing.T) {

		// given
		distances := map[crossCountryTravel]Kilometers{
			crossCountryTravel{from: "A", to: "B"}: 100,
		}
		speed := KilometersPerHour(10)

		estimatorUT := NewSpeedBasedTravelsEstimator(distances, speed)

		// when
		estimation, err := estimatorUT.EstimateTravelDuration("B", "C")

		// then
		require.Error(t, err)
		assert.Empty(t, estimation)
	})

	t.Run(fmt.Sprintf("with speed of light"), func(t *testing.T) {
		estimatorUT := NewSpeedOfLightTravelsEstimator()
		for travel := range sampleDistances {
			t.Run(fmt.Sprintf("travel %v should last 0", travel), func(t *testing.T) {
				// when
				estimation, err := estimatorUT.EstimateTravelDuration(travel.from, travel.to)
				// then
				require.NoError(t, err)
				assert.Empty(t, estimation)
			})
		}
	})

	t.Run(fmt.Sprintf("with avg speed of plane"), func(t *testing.T) {
		estimatorUT := NewAvgAirplaneTravelsEstimator()
		for travel := range sampleDistances {
			t.Run(fmt.Sprintf("travel %v should last more than 0", travel), func(t *testing.T) {
				// when
				estimation, err := estimatorUT.EstimateTravelDuration(travel.from, travel.to)
				// then
				require.NoError(t, err)
				assert.True(t, estimation > time.Duration(0))
			})
		}
	})
}
