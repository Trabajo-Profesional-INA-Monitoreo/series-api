package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
)

type StreamService interface {
	GetNetworks() dtos.StreamsPerNetworkResponse
	GetStations()
}

type streamService struct {
	repository repositories.StreamRepository
}

func NewSeriesService(repository repositories.StreamRepository) StreamService {
	return &streamService{repository}
}

func (s streamService) GetNetworks() dtos.StreamsPerNetworkResponse {
	networks := s.repository.GetNetworks()
	return dtos.StreamsPerNetworkResponse{Networks: networks}
}

func (s streamService) GetStations() {

}
