package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"time"
)

type InaServiceApi interface {
	GetObservatedSerieById(id uint64, start time.Time, end time.Time) (dtos.StreamsDataResponse, error)
	GetPredictedSerieById(id uint64, streamId uint64) (dtos.CalibratedStreamsDataResponse, error)
}

type inaServiceApiImpl struct {
	inaApiClient clients.InaAPiClient
}

func NewInaServiceApi(inaApiClient clients.InaAPiClient) InaServiceApi {
	return &inaServiceApiImpl{inaApiClient: inaApiClient}
}

func (s inaServiceApiImpl) GetObservatedSerieById(id uint64, start time.Time, end time.Time) (dtos.StreamsDataResponse, error) {
	return s.getDataFromObservedStream(id, start, end)
}

func (s inaServiceApiImpl) getDataFromObservedStream(id uint64, start time.Time, end time.Time) (dtos.StreamsDataResponse, error) {
	streams, err := s.inaApiClient.GetObservedData(id, start, end)
	if err != nil {
		return dtos.StreamsDataResponse{}, err
	}
	var streamsData []dtos.StreamsData
	for _, stream := range streams {
		streamsData = append(streamsData, stream.ConvertToStreamData())
	}
	return dtos.StreamsDataResponse{Streams: streamsData}, nil
}

func (s inaServiceApiImpl) GetPredictedSerieById(id uint64, streamId uint64) (dtos.CalibratedStreamsDataResponse, error) {
	streams, err := s.inaApiClient.GetLastForecast(id)
	if err != nil {
		return dtos.CalibratedStreamsDataResponse{}, err
	}
	return streams.ConvertToCalibratedStreamsDataResponse(streamId), nil
}
