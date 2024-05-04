package inputs_service

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	"time"
)

type InputsService interface {
	GetGeneralMetrics(configurationId uint64) dtos.InputsGeneralMetrics
	GetTotalStreamsWithNullValues(configurationId uint64, timeStart time.Time, timeEnd time.Time) dtos.TotalStreamsWithNullValues
	GetTotalStreamsWithObservedOutlier(configurationId uint64, timeStart time.Time, timeEnd time.Time) dtos.TotalStreamsWithObservedOutlier
}

type inputsService struct {
	repository repositories.InputsRepository
}

func NewInputsService(repository repositories.InputsRepository) InputsService {
	return &inputsService{repository}
}

func (s inputsService) GetGeneralMetrics(configurationId uint64) dtos.InputsGeneralMetrics {
	streamsResult := make(chan int, 1)
	stationsResult := make(chan int, 1)
	go func() {
		streamsResult <- s.repository.GetTotalStreams(configurationId)
	}()
	go func() {
		stationsResult <- s.repository.GetTotalStations(configurationId)
	}()
	totalStreams := <-streamsResult
	totalStations := <-stationsResult
	return dtos.InputsGeneralMetrics{TotalStreams: totalStreams, TotalStations: totalStations}
}

func (s inputsService) GetTotalStreamsWithNullValues(configurationId uint64, timeStart time.Time, timeEnd time.Time) dtos.TotalStreamsWithNullValues {
	streamsResult := make(chan int, 1)
	go func() {
		streamsResult <- s.repository.GetTotalStreams(configurationId)
	}()
	streamsWithNull := s.repository.GetTotalStreamsByError(configurationId, timeStart, timeEnd, entities.NullValue)
	totalStreams := <-streamsResult

	return dtos.TotalStreamsWithNullValues{TotalStreams: totalStreams, TotalStreamsWithNull: streamsWithNull}
}

func (s inputsService) GetTotalStreamsWithObservedOutlier(configurationId uint64, timeStart time.Time, timeEnd time.Time) dtos.TotalStreamsWithObservedOutlier {
	streamsResult := make(chan int, 1)
	go func() {
		streamsResult <- s.repository.GetTotalStreams(configurationId)
	}()
	streamsWithObservedOutlier := s.repository.GetTotalStreamsByError(configurationId, timeStart, timeEnd, entities.ObservedOutlier)
	totalStreams := <-streamsResult

	return dtos.TotalStreamsWithObservedOutlier{TotalStreams: totalStreams, TotalStreamsWithObservedOutlier: streamsWithObservedOutlier}
}
