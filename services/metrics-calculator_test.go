package services

import (
	"errors"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients/responses"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type MockInaApiClient struct {
	ObservedData []responses.ObservedDataResponse
	Error        error
	LastForecast responses.LastForecast
}

func (m MockInaApiClient) GetLastForecast(calibrationId uint64) (*responses.LastForecast, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return &m.LastForecast, nil
}

func (m MockInaApiClient) GetObservedData(streamId uint64, timeStart time.Time, timeEnd time.Time) ([]responses.ObservedDataResponse, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.ObservedData, nil
}

func (m MockInaApiClient) GetStream(streamId uint64) (*responses.InaStreamResponse, error) {
	return nil, nil
}

func TestShouldReturnLevelCountForAStream(t *testing.T) {

	mockInaApiClient := MockInaApiClient{}
	evacuationLevel := 5.0
	alertLevel := 3.0
	lowWaterLevel := 1.0
	streams := []dtos.BehaviourStream{
		{
			StreamId:        1,
			EvacuationLevel: &evacuationLevel,
			AlertLevel:      &alertLevel,
			LowWaterLevel:   &lowWaterLevel,
		},
	}

	var observedData []responses.ObservedDataResponse

	for i := 0; i < 10; i++ {
		level := float64(i)
		observedData = append(observedData, responses.ObservedDataResponse{Value: &level})
	}

	mockInaApiClient.ObservedData = observedData

	result := getLevelsCountForAllStreams(streams, time.Now(), time.Now(), mockInaApiClient)

	assert.Equal(t, uint64(10), result.TotalValuesCount)
	assert.Equal(t, uint64(5), result.CountEvacuationLevel)
	assert.Equal(t, uint64(2), result.CountAlertLevel)
	assert.Equal(t, uint64(2), result.CountLowWaterLevel)
}

func TestShouldReturnAnEmptyListIfThereIsAnErrorFetchingTheData(t *testing.T) {

	mockInaApiClient := MockInaApiClient{}
	evacuationLevel := 5.0
	alertLevel := 3.0
	lowWaterLevel := 1.0
	streams := []dtos.BehaviourStream{
		{
			StreamId:        1,
			EvacuationLevel: &evacuationLevel,
			AlertLevel:      &alertLevel,
			LowWaterLevel:   &lowWaterLevel,
		},
	}

	mockInaApiClient.Error = errors.New("error fetching data")

	result := getLevelsCountForAllStreams(streams, time.Now(), time.Now(), mockInaApiClient)

	assert.Equal(t, uint64(0), result.TotalValuesCount)
	assert.Equal(t, uint64(0), result.CountEvacuationLevel)
	assert.Equal(t, uint64(0), result.CountAlertLevel)
	assert.Equal(t, uint64(0), result.CountLowWaterLevel)
}
