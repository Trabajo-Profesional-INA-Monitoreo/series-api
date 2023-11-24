package services

import "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"

type StreamService interface {
	GetNetworks()
	GetStations()
}

type streamService struct {
	repository repositories.StreamRepository
}

func NewSeriesService(repository repositories.StreamRepository) StreamService {
	return &streamService{repository}
}

func (s streamService) GetNetworks() {

}

func (s streamService) GetStations() {

}
