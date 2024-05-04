package services

import (
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

func TestShouldReturnTheStationsSummary(t *testing.T) {
	mockRepository := stationsRepositoryMock{}
	mockRepository.streamsPerStation = []*dtos.StreamsPerStation{{"Test", "1", 1, 0}}
	mockRepository.errorsPerStation = []dtos.ErrorsOfStations{{"1", 1}}
	stationsService := &stationsServiceImpl{mockRepository}
	res := stationsService.GetStations(time.Now(), time.Now(), 1)
	assert.Equal(t, 1, len(res.Stations))
	assert.Equal(t, 1, res.Stations[0].StreamsCount)
	assert.Equal(t, 1, res.Stations[0].ErrorCount)
	assert.Equal(t, "Test", res.Stations[0].StationName)
	assert.Equal(t, "1", res.Stations[0].StationId)
}
