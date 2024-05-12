package stations_service

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	metrics_service "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services/metrics-service"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type StationsService interface {
	GetStations(time.Time, time.Time, uint64) dtos.StreamsPerStationResponse
}

type stationsServiceImpl struct {
	repository repositories.StationsRepository
	inaClient  clients.InaAPiClient
}

func NewStationsService(repository repositories.StationsRepository, inaClient clients.InaAPiClient) StationsService {
	return &stationsServiceImpl{repository: repository, inaClient: inaClient}
}

func (s stationsServiceImpl) GetStations(timeStart time.Time, timeEnd time.Time, configId uint64) dtos.StreamsPerStationResponse {
	stations := s.repository.GetStations(configId)
	if stations == nil {
		return dtos.StreamsPerStationResponse{}
	}
	for _, station := range *stations {
		id, _ := strconv.ParseUint(station.StationId, 10, 64)
		mainStream, err := s.inaClient.GetMainStreamFromStation(id)
		if err != nil {
			log.Errorf("Error getting station streams: %v for stations summary", err)
			continue
		}
		if mainStream == nil {
			continue
		}
		station.MainStreamId = &mainStream.StreamId
		station.LastUpdate = mainStream.LastUpdate
		levels, err := s.inaClient.GetObservedData(mainStream.StreamId, timeStart, timeEnd)
		if err != nil {
			log.Errorf("Error getting levels: %v for stations summary", err)
			continue
		}
		calculator := metrics_service.NewCalculatorOfWaterLevels(mainStream.AlertLevel, nil, nil)
		totalValues := uint64(0)
		for _, level := range levels {
			if level.Value != nil {
				calculator.Compute(*level.Value)
				totalValues++
			}
		}
		station.AlertWaterLevels = calculator.GetAlertsCount()
		station.TotalWaterLevels = totalValues
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
