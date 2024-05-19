package stations_service

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	metrics_service "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services/metrics-service"
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
)

type StationsService interface {
	GetStations(*dtos.QueryParameters) dtos.StreamsPerStationResponse
}

type stationsServiceImpl struct {
	repository repositories.StationsRepository
	inaClient  clients.InaAPiClient
}

func NewStationsService(repository repositories.StationsRepository, inaClient clients.InaAPiClient) StationsService {
	return &stationsServiceImpl{repository: repository, inaClient: inaClient}
}

func (s stationsServiceImpl) GetStations(parameters *dtos.QueryParameters) dtos.StreamsPerStationResponse {
	configId := parameters.Get("configurationId").(uint64)
	page := *parameters.GetAsInt("page")
	pageSize := *parameters.GetAsInt("pageSize")
	timeStart := parameters.Get("timeStart").(time.Time)
	timeEnd := parameters.Get("timeEnd").(time.Time)
	stations, pageable := s.repository.GetStations(configId, page, pageSize)
	if stations == nil {
		return dtos.StreamsPerStationResponse{}
	}
	var wg sync.WaitGroup
	for _, station := range *stations {
		wg.Add(1)
		go s.getLastUpdateAndLevelFromStation(station, timeStart, timeEnd, &wg)
	}

	var stationIds []uint64
	for _, station := range *stations {
		id, _ := strconv.ParseUint(station.StationId, 10, 64)
		stationIds = append(stationIds, id)
	}
	errorsPerStation := s.repository.GetErrorsOfStations(configId, timeStart, timeEnd, stationIds)

	for _, errors := range errorsPerStation {
		for _, station := range *stations {
			if station.StationId == errors.StationId {
				station.ErrorCount = errors.ErrorCount
				break
			}
		}
	}
	wg.Wait()
	return dtos.StreamsPerStationResponse{Stations: *stations, Pageable: pageable}
}

func (s stationsServiceImpl) getLastUpdateAndLevelFromStation(station *dtos.StreamsPerStation, timeStart time.Time, timeEnd time.Time, wg *sync.WaitGroup) {
	defer wg.Done()
	id, _ := strconv.ParseUint(station.StationId, 10, 64)
	mainStream, err := s.inaClient.GetMainStreamFromStation(id)
	if err != nil {
		log.Errorf("Error getting station streams: %v for stations summary", err)
		return
	}
	if mainStream == nil {
		return
	}
	station.MainStreamId = &mainStream.StreamId
	station.LastUpdate = mainStream.LastUpdate
	levels, err := s.inaClient.GetObservedData(mainStream.StreamId, timeStart, timeEnd)
	if err != nil {
		log.Errorf("Error getting levels: %v for stations summary", err)
		return
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
