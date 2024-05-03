package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type inputsRepositoryMock struct {
	totalStreams  int
	totalStations int
	totalErrors   int
}

func (i inputsRepositoryMock) GetTotalStreams(configurationId uint64) int {
	return i.totalStreams
}

func (i inputsRepositoryMock) GetTotalStations(configurationId uint64) int {
	return i.totalStations
}

func (i inputsRepositoryMock) GetTotalStreamsByError(id uint64, start time.Time, end time.Time, value entities.ErrorType) int {
	return i.totalErrors
}

func TestShouldReturnGeneralInputMetrics(t *testing.T) {
	mockRepository := inputsRepositoryMock{1, 1, 0}
	inputsService := &inputsService{mockRepository}
	metrics := inputsService.GetGeneralMetrics(1)
	assert.Equal(t, 1, metrics.TotalStations)
	assert.Equal(t, 1, metrics.TotalStreams)
}

func TestShouldReturnNullValuesInput(t *testing.T) {
	mockRepository := inputsRepositoryMock{5, 1, 1}
	inputsService := &inputsService{mockRepository}
	metrics := inputsService.GetTotalStreamsWithNullValues(1, time.Now(), time.Now())
	assert.Equal(t, 1, metrics.TotalStreamsWithNull)
	assert.Equal(t, 5, metrics.TotalStreams)
}

func TestShouldReturnOutlierValuesInput(t *testing.T) {
	mockRepository := inputsRepositoryMock{5, 1, 1}
	inputsService := &inputsService{mockRepository}
	metrics := inputsService.GetTotalStreamsWithObservedOutlier(1, time.Now(), time.Now())
	assert.Equal(t, 1, metrics.TotalStreamsWithObservedOutlier)
	assert.Equal(t, 5, metrics.TotalStreams)
}
