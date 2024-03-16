package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	"strconv"
	"time"
)

type StreamService interface {
	GetNetworks() dtos.StreamsPerNetworkResponse
	GetStations() dtos.StreamsPerStationResponse
	GetCuredSerieById(id string) dtos.StreamsDataResponse
	GetObservatedSerieById(id string) dtos.StreamsDataResponse
	GetPredictedSerieById(id string) dtos.CalibratedStreamsDataResponse
}

type streamService struct {
	repository   repositories.StreamRepository
	inaApiClient clients.InaAPiClient
}

func NewSeriesService(repository repositories.StreamRepository, inaApiClient clients.InaAPiClient) StreamService {
	return &streamService{repository, inaApiClient}
}

func (s streamService) GetNetworks() dtos.StreamsPerNetworkResponse {
	networks := s.repository.GetNetworks()
	return dtos.StreamsPerNetworkResponse{Networks: networks}
}

func (s streamService) GetStations() dtos.StreamsPerStationResponse {
	stations := s.repository.GetStations()
	return dtos.StreamsPerStationResponse{Stations: stations}
}

func (s streamService) GetCuredSerieById(id string) dtos.StreamsDataResponse {
	today := time.Now()
	num, _ := strconv.ParseUint(id, 10, 64)
	streams, _ := s.inaApiClient.GetObservedData(num, today.Add(time.Duration(-24*5)*time.Hour), today.Add(time.Duration(24*5)*time.Hour))
	var streamsData []dtos.StreamsData
	for _, stream := range streams {
		streamsData = append(streamsData, stream.ConvertToStreamData())
	}
	return dtos.StreamsDataResponse{Streams: streamsData}
}

func (s streamService) GetObservatedSerieById(id string) dtos.StreamsDataResponse {
	today := time.Now()
	num, _ := strconv.ParseUint(id, 10, 64)
	streams, _ := s.inaApiClient.GetObservedData(num, today.Add(time.Duration(-24*5)*time.Hour), today.Add(time.Duration(24*1)*time.Hour))
	var streamsData []dtos.StreamsData
	for _, stream := range streams {
		streamsData = append(streamsData, stream.ConvertToStreamData())
	}
	return dtos.StreamsDataResponse{Streams: streamsData}
}

func (s streamService) GetPredictedSerieById(id string) dtos.CalibratedStreamsDataResponse {
	num, _ := strconv.ParseUint(id, 10, 64)
	streams, _ := s.inaApiClient.GetLastForecast(num)
	return streams.ConvertToCalibratedStreamsDataResponse()
}
