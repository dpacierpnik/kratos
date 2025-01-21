package detector

import (
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDetector(t *testing.T) {

	t.Run("should return an error if fetching aggregate fails", func(t *testing.T) {

		// given
		visitsAggregator := &visitsAggregatorMock{}
		travelDurationEstimator := &travelDurationEstimatorMock{}

		previousVisitsAggregate := &previousVisitsAggregateMock{}

		identityID := uuid.Must(uuid.NewV4())
		currentVisit := Visit[anythingGeoLocated]{}

		visitsAggregator.ExpectGetPreviousVisitsAggregate(identityID).Return(previousVisitsAggregate, errors.New("test error"))
		visitsAggregator.ExpectAggregateVisit(identityID, currentVisit, DetectionError).Return(nil)

		detectorUT := NewDetector(visitsAggregator, travelDurationEstimator)

		// when
		result := detectorUT.DetectImpossibleTravel(identityID, currentVisit)

		// then
		require.Equal(t, DetectionError, result)

		visitsAggregator.AssertExpectations(t)
		travelDurationEstimator.AssertExpectations(t)
	})

	t.Run("should detect legit first visit", func(t *testing.T) {

		// given
		visitsAggregator := &visitsAggregatorMock{}
		travelDurationEstimator := &travelDurationEstimatorMock{}

		previousVisitsAggregate := &previousVisitsAggregateMock{}
		previousVisitsAggregate.ExpectGetLastLegitVisit().Return(Visit[anythingGeoLocated]{}, false)

		identityID := uuid.Must(uuid.NewV4())
		currentVisit := Visit[anythingGeoLocated]{}

		visitsAggregator.ExpectGetPreviousVisitsAggregate(identityID).Return(previousVisitsAggregate, nil)
		visitsAggregator.ExpectAggregateVisit(identityID, currentVisit, LegitFirstVisit).Return(nil)

		detectorUT := NewDetector(visitsAggregator, travelDurationEstimator)

		// when
		result := detectorUT.DetectImpossibleTravel(identityID, currentVisit)

		// then
		require.Equal(t, LegitFirstVisit, result)

		visitsAggregator.AssertExpectations(t)
		travelDurationEstimator.AssertExpectations(t)
	})

	t.Run("should detect legit visit by preconditions", func(t *testing.T) {

		// given
		visitsAggregator := &visitsAggregatorMock{}
		travelDurationEstimator := &travelDurationEstimatorMock{}
		previousVisitsAggregate := &previousVisitsAggregateMock{}

		identityID := uuid.Must(uuid.NewV4())
		currentVisit := Visit[anythingGeoLocated]{}

		previousVisitsAggregate.ExpectGetLastLegitVisit().Return(Visit[anythingGeoLocated]{}, true)
		previousVisitsAggregate.ExpectIsLegitVisit(currentVisit).Return(true)

		visitsAggregator.ExpectGetPreviousVisitsAggregate(identityID).Return(previousVisitsAggregate, nil)
		visitsAggregator.ExpectAggregateVisit(identityID, currentVisit, LegitVisit).Return(nil)

		detectorUT := NewDetector(visitsAggregator, travelDurationEstimator)

		// when
		result := detectorUT.DetectImpossibleTravel(identityID, currentVisit)

		// then
		require.Equal(t, LegitVisit, result)

		visitsAggregator.AssertExpectations(t)
		travelDurationEstimator.AssertExpectations(t)
	})

	t.Run("should not be able to detect impossible travel if last visit is not geo-located", func(t *testing.T) {

		// given
		visitsAggregator := &visitsAggregatorMock{}
		travelDurationEstimator := &travelDurationEstimatorMock{}
		previousVisitsAggregate := &previousVisitsAggregateMock{}

		identityID := uuid.Must(uuid.NewV4())
		deviceID := uuid.Must(uuid.NewV4())
		currentVisit := Visit[anythingGeoLocated]{
			Time:        time.Now(),
			DeviceId:    deviceID,
			GeoLocation: "location_1",
		}
		lastLegitVisit := Visit[anythingGeoLocated]{
			Time:        time.Now(),
			DeviceId:    deviceID,
			GeoLocation: "",
		}

		previousVisitsAggregate.ExpectGetLastLegitVisit().Return(lastLegitVisit, true)
		previousVisitsAggregate.ExpectIsLegitVisit(currentVisit).Return(false)

		visitsAggregator.ExpectGetPreviousVisitsAggregate(identityID).Return(previousVisitsAggregate, nil)
		visitsAggregator.ExpectAggregateVisit(identityID, currentVisit, InsufficientData).Return(nil)

		detectorUT := NewDetector(visitsAggregator, travelDurationEstimator)

		// when
		result := detectorUT.DetectImpossibleTravel(identityID, currentVisit)

		// then
		require.Equal(t, InsufficientData, result)

		visitsAggregator.AssertExpectations(t)
		travelDurationEstimator.AssertExpectations(t)
	})

	t.Run("should not be able to detect impossible travel if current visit is not geo-located", func(t *testing.T) {

		// given
		visitsAggregator := &visitsAggregatorMock{}
		travelDurationEstimator := &travelDurationEstimatorMock{}
		previousVisitsAggregate := &previousVisitsAggregateMock{}

		identityID := uuid.Must(uuid.NewV4())
		deviceID := uuid.Must(uuid.NewV4())
		currentVisit := Visit[anythingGeoLocated]{
			Time:        time.Now(),
			DeviceId:    deviceID,
			GeoLocation: "",
		}
		lastLegitVisit := Visit[anythingGeoLocated]{
			Time:        time.Now(),
			DeviceId:    deviceID,
			GeoLocation: "location_2",
		}

		previousVisitsAggregate.ExpectGetLastLegitVisit().Return(lastLegitVisit, true)
		previousVisitsAggregate.ExpectIsLegitVisit(currentVisit).Return(false)

		visitsAggregator.ExpectGetPreviousVisitsAggregate(identityID).Return(previousVisitsAggregate, nil)
		visitsAggregator.ExpectAggregateVisit(identityID, currentVisit, InsufficientData).Return(nil)

		detectorUT := NewDetector(visitsAggregator, travelDurationEstimator)

		// when
		result := detectorUT.DetectImpossibleTravel(identityID, currentVisit)

		// then
		require.Equal(t, InsufficientData, result)

		visitsAggregator.AssertExpectations(t)
		travelDurationEstimator.AssertExpectations(t)
	})

	t.Run("should detect impossible travel for unreasonable short period of time", func(t *testing.T) {

		// given
		visitsAggregator := &visitsAggregatorMock{}
		travelDurationEstimator := &travelDurationEstimatorMock{}
		previousVisitsAggregate := &previousVisitsAggregateMock{}

		identityID := uuid.Must(uuid.NewV4())
		deviceID := uuid.Must(uuid.NewV4())
		currentVisit := Visit[anythingGeoLocated]{
			Time:        time.Now(),
			DeviceId:    deviceID,
			GeoLocation: "location_1",
		}
		lastLegitVisit := Visit[anythingGeoLocated]{
			Time:        time.Now().Add(1 * time.Minute),
			DeviceId:    deviceID,
			GeoLocation: "location_2",
		}

		previousVisitsAggregate.ExpectGetLastLegitVisit().Return(lastLegitVisit, true)
		previousVisitsAggregate.ExpectIsLegitVisit(currentVisit).Return(false)

		visitsAggregator.ExpectGetPreviousVisitsAggregate(identityID).Return(previousVisitsAggregate, nil)
		visitsAggregator.ExpectAggregateVisit(identityID, currentVisit, ImpossibleTravel).Return(nil)

		travelDurationEstimator.ExpectEstimateTravelDuration("location_2", "location_1").Return(time.Minute*10, nil)

		detectorUT := NewDetector(visitsAggregator, travelDurationEstimator)

		// when
		result := detectorUT.DetectImpossibleTravel(identityID, currentVisit)

		// then
		require.Equal(t, ImpossibleTravel, result)

		visitsAggregator.AssertExpectations(t)
		travelDurationEstimator.AssertExpectations(t)
	})

	t.Run("should return detection error if travel duration estimation fails", func(t *testing.T) {

		// given
		visitsAggregator := &visitsAggregatorMock{}
		travelDurationEstimator := &travelDurationEstimatorMock{}
		previousVisitsAggregate := &previousVisitsAggregateMock{}

		identityID := uuid.Must(uuid.NewV4())
		deviceID := uuid.Must(uuid.NewV4())
		currentVisit := Visit[anythingGeoLocated]{
			Time:        time.Now(),
			DeviceId:    deviceID,
			GeoLocation: "location_1",
		}
		lastLegitVisit := Visit[anythingGeoLocated]{
			Time:        time.Now().Add(1 * time.Minute),
			DeviceId:    deviceID,
			GeoLocation: "location_2",
		}

		previousVisitsAggregate.ExpectGetLastLegitVisit().Return(lastLegitVisit, true)
		previousVisitsAggregate.ExpectIsLegitVisit(currentVisit).Return(false)

		visitsAggregator.ExpectGetPreviousVisitsAggregate(identityID).Return(previousVisitsAggregate, nil)
		visitsAggregator.ExpectAggregateVisit(identityID, currentVisit, DetectionError).Return(nil)

		travelDurationEstimator.ExpectEstimateTravelDuration("location_2", "location_1").Return(0, errors.New("test error"))

		detectorUT := NewDetector(visitsAggregator, travelDurationEstimator)

		// when
		result := detectorUT.DetectImpossibleTravel(identityID, currentVisit)

		// then
		require.Equal(t, DetectionError, result)

		visitsAggregator.AssertExpectations(t)
		travelDurationEstimator.AssertExpectations(t)
	})

	t.Run("should detect legit travel for reasonable period of time", func(t *testing.T) {

		// given
		visitsAggregator := &visitsAggregatorMock{}
		travelDurationEstimator := &travelDurationEstimatorMock{}
		previousVisitsAggregate := &previousVisitsAggregateMock{}

		identityID := uuid.Must(uuid.NewV4())
		deviceID := uuid.Must(uuid.NewV4())
		currentVisit := Visit[anythingGeoLocated]{
			Time:        time.Now(),
			DeviceId:    deviceID,
			GeoLocation: "location_1",
		}
		lastLegitVisit := Visit[anythingGeoLocated]{
			Time:        time.Now().Add(2 * time.Minute),
			DeviceId:    deviceID,
			GeoLocation: "location_2",
		}

		previousVisitsAggregate.ExpectGetLastLegitVisit().Return(lastLegitVisit, true)
		previousVisitsAggregate.ExpectIsLegitVisit(currentVisit).Return(false)

		visitsAggregator.ExpectGetPreviousVisitsAggregate(identityID).Return(previousVisitsAggregate, nil)
		visitsAggregator.ExpectAggregateVisit(identityID, currentVisit, LegitTravel).Return(nil)

		travelDurationEstimator.ExpectEstimateTravelDuration("location_2", "location_1").Return(time.Minute*1, nil)

		detectorUT := NewDetector(visitsAggregator, travelDurationEstimator)

		// when
		result := detectorUT.DetectImpossibleTravel(identityID, currentVisit)

		// then
		require.Equal(t, LegitTravel, result)

		visitsAggregator.AssertExpectations(t)
		travelDurationEstimator.AssertExpectations(t)
	})

	t.Run("should detect legit travel even if aggregation of current visit fails", func(t *testing.T) {

		// given
		visitsAggregator := &visitsAggregatorMock{}
		travelDurationEstimator := &travelDurationEstimatorMock{}
		previousVisitsAggregate := &previousVisitsAggregateMock{}

		identityID := uuid.Must(uuid.NewV4())
		deviceID := uuid.Must(uuid.NewV4())
		currentVisit := Visit[anythingGeoLocated]{
			Time:        time.Now(),
			DeviceId:    deviceID,
			GeoLocation: "location_1",
		}
		lastLegitVisit := Visit[anythingGeoLocated]{
			Time:        time.Now().Add(2 * time.Minute),
			DeviceId:    deviceID,
			GeoLocation: "location_2",
		}

		previousVisitsAggregate.ExpectGetLastLegitVisit().Return(lastLegitVisit, true)
		previousVisitsAggregate.ExpectIsLegitVisit(currentVisit).Return(false)

		visitsAggregator.ExpectGetPreviousVisitsAggregate(identityID).Return(previousVisitsAggregate, nil)
		visitsAggregator.ExpectAggregateVisit(identityID, currentVisit, LegitTravel).Return(errors.New("test error"))

		travelDurationEstimator.ExpectEstimateTravelDuration("location_2", "location_1").Return(time.Minute*1, nil)

		detectorUT := NewDetector(visitsAggregator, travelDurationEstimator)

		// when
		result := detectorUT.DetectImpossibleTravel(identityID, currentVisit)

		// then
		require.Equal(t, LegitTravel, result)

		visitsAggregator.AssertExpectations(t)
		travelDurationEstimator.AssertExpectations(t)
	})
}

type anythingGeoLocated string

func (l anythingGeoLocated) IsGeoLocated() bool {
	return l != ""
}

type previousVisitsAggregateMock struct {
	mock.Mock
}

func (m *previousVisitsAggregateMock) GetLastLegitVisit() (Visit[anythingGeoLocated], bool) {
	args := m.Called()
	v, _ := args.Get(0).(Visit[anythingGeoLocated])
	return v, args.Bool(1)
}

func (m *previousVisitsAggregateMock) IsLegitVisit(currentVisit Visit[anythingGeoLocated]) bool {
	args := m.Called(currentVisit)
	return args.Bool(0)
}

func (m *previousVisitsAggregateMock) ExpectGetLastLegitVisit() *mock.Call {
	return m.On("GetLastLegitVisit")
}

func (m *previousVisitsAggregateMock) ExpectIsLegitVisit(currentVisit Visit[anythingGeoLocated]) *mock.Call {
	return m.On("IsLegitVisit", currentVisit)
}

type visitsAggregatorMock struct {
	mock.Mock
}

func (m *visitsAggregatorMock) GetPreviousVisitsAggregate(identityID uuid.UUID) (PreviousVisitsAggregate[anythingGeoLocated], error) {
	args := m.Called(identityID)
	agg, _ := args.Get(0).(PreviousVisitsAggregate[anythingGeoLocated])
	return agg, args.Error(1)
}

func (m *visitsAggregatorMock) AggregateVisit(identityID uuid.UUID, visit Visit[anythingGeoLocated], detectionResult DetectionResult) error {
	args := m.Called(identityID, visit, detectionResult)
	return args.Error(0)
}

func (m *visitsAggregatorMock) ExpectGetPreviousVisitsAggregate(identityID uuid.UUID) *mock.Call {
	return m.On("GetPreviousVisitsAggregate", identityID)
}

func (m *visitsAggregatorMock) ExpectAggregateVisit(identityID uuid.UUID, visit Visit[anythingGeoLocated], detectionResult DetectionResult) *mock.Call {
	return m.On("AggregateVisit", identityID, visit, detectionResult)
}

type travelDurationEstimatorMock struct {
	mock.Mock
}

func (m *travelDurationEstimatorMock) EstimateTravelDuration(from anythingGeoLocated, to anythingGeoLocated) (time.Duration, error) {
	args := m.Called(from, to)
	d, _ := args.Get(0).(time.Duration)
	return d, args.Error(1)
}

func (m *travelDurationEstimatorMock) ExpectEstimateTravelDuration(from anythingGeoLocated, to anythingGeoLocated) *mock.Call {
	return m.On("EstimateTravelDuration", from, to)
}
