package stations_service

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	"time"
)

type StationsService interface {
	GetStations(time.Time, time.Time, uint64) dtos.StreamsPerStationResponse
}

type stationsServiceImpl struct {
	repository repositories.StationsRepository
}

func NewStationsService(repository repositories.StationsRepository) StationsService {
	return &stationsServiceImpl{repository: repository}
}

func (s stationsServiceImpl) GetStations(timeStart time.Time, timeEnd time.Time, configId uint64) dtos.StreamsPerStationResponse {
	stations := s.repository.GetStations(configId)
	if stations == nil {
		return dtos.StreamsPerStationResponse{}
	}
	errorsPerStation := s.repository.GetErrorsOfStations(configId, timeStart, timeEnd)

	for _, errors := range errorsPerStation {
		for _, station := range *stations {
			if station.StationId == errors.StationId {
				station.ErrorCount = errors.ErrorCount
				break
			}
		}
	}
	return dtos.StreamsPerStationResponse{Stations: *stations}
}