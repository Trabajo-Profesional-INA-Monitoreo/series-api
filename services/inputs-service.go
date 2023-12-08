package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
)

type InputsService interface {
	GetGeneralMetrics() dtos.InputsGeneralMetrics
}

type inputsService struct {
	repository repositories.StreamRepository
}

func NewInputsService(repository repositories.StreamRepository) InputsService {
	return &inputsService{repository}
}

func (s inputsService) GetGeneralMetrics() dtos.InputsGeneralMetrics {
	streamsResult := make(chan int, 1)
	stationsResult := make(chan int, 1)
	go func() {
		streamsResult <- s.repository.GetTotalStreams()
	}()
	go func() {
		stationsResult <- s.repository.GetTotalStations()
	}()
	totalNetworks := s.repository.GetTotalNetworks()
	totalStreams := <-streamsResult
	totalStations := <-stationsResult
	return dtos.InputsGeneralMetrics{TotalStreams: totalStreams, TotalStations: totalStations, TotalNetworks: totalNetworks}
}
