package stations_service

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients/responses"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type stationsRepositoryMock struct {
	streamsPerStation []*dtos.StreamsPerStation
	errorsPerStation  []dtos.ErrorsOfStations
}

func (s stationsRepositoryMock) GetStations(configId uint64) *[]*dtos.StreamsPerStation {
	return &s.streamsPerStation

}
func (s stationsRepositoryMock) GetErrorsOfStations(configId uint64, timeStart time.Time, timeEnd time.Time) []dtos.ErrorsOfStations {
	return s.errorsPerStation
}

type MockInaApiClient struct {
	ObservedData []responses.ObservedDataResponse
	Error        error
	Stream       responses.InaStreamResponse
}

func (m MockInaApiClient) GetLastForecast(calibrationId uint64) (*responses.LastForecast, error) {
	return nil, nil
}

func (m MockInaApiClient) GetObservedData(streamId uint64, timeStart time.Time, timeEnd time.Time) ([]responses.ObservedDataResponse, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.ObservedData, nil
}

func (m MockInaApiClient) GetStream(streamId uint64) (*responses.InaStreamResponse, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return &m.Stream, nil
}

func (m MockInaApiClient) GetMainStreamFromStation(stationId uint64) (*responses.MainStreamFromStation, error) {
	return responses.NewMainStreamFromStation(m.Stream), nil
}

func TestShouldReturnTheStationsSummary(t *testing.T) {
	mockRepository := stationsRepositoryMock{}
	mockRepository.streamsPerStation = []*dtos.StreamsPerStation{{"Test", "1", 1, 0, nil, 0, 0, nil}}
	mockRepository.errorsPerStation = []dtos.ErrorsOfStations{{"1", 1}}
	mockClient := MockInaApiClient{}
	alertLevel := 4.0
	mockClient.ObservedData = []responses.ObservedDataResponse{{Value: &alertLevel}}
	stationsService := &stationsServiceImpl{mockRepository, mockClient}
	res := stationsService.GetStations(time.Now(), time.Now(), 1)
	assert.Equal(t, 1, len(res.Stations))
	assert.Equal(t, 1, res.Stations[0].StreamsCount)
	assert.Equal(t, 1, res.Stations[0].ErrorCount)
	assert.Equal(t, "Test", res.Stations[0].StationName)
	assert.Equal(t, "1", res.Stations[0].StationId)
}
