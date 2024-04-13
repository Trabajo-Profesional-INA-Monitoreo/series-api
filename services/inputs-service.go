package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	"time"
)

type InputsService interface {
	GetGeneralMetrics(configurationId uint64) dtos.InputsGeneralMetrics
	GetTotalStreamsWithNullValues(configurationId uint64, timeStart time.Time, timeEnd time.Time) dtos.TotalStreamsWithNullValues
}

type inputsService struct {
	repository repositories.StreamRepository
}

func NewInputsService(repository repositories.StreamRepository) InputsService {
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
	streamsWithNull := s.repository.GetTotalStreamsWithNullValues(configurationId, timeStart, timeEnd)
	totalStreams := <-streamsResult

	return dtos.TotalStreamsWithNullValues{TotalStreams: totalStreams, TotalStreamsWithNull: streamsWithNull}
}
