package outputs_service

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services/metrics-service"
	log "github.com/sirupsen/logrus"
	"time"
)

type OutputsService interface {
	GetOutputBehaviourMetrics(configId uint64, timeStart time.Time, timeEnd time.Time) (*dtos.BehaviourStreamsResponse, error)
}

type outputsServiceImpl struct {
	inaApiClient clients.InaAPiClient
	repository   repositories.OutputsRepository
}

func NewOutputsService(repository repositories.OutputsRepository, inaApiClient clients.InaAPiClient) OutputsService {
	return &outputsServiceImpl{repository: repository, inaApiClient: inaApiClient}
}

func (s outputsServiceImpl) getLevelsCountForAllStreams(behaviourStreams []dtos.BehaviourStream, timeStart time.Time, timeEnd time.Time) *dtos.BehaviourStreamsResponse {
	behaviourResponse := dtos.NewBehaviourStreamsResponse()
	for _, stream := range behaviourStreams {
		values, err := s.inaApiClient.GetObservedData(stream.StreamId, timeStart, timeEnd)
		if err != nil {
			log.Errorf("GetOutputBehaviourMetrics | Could not get metrics for stream with id %v: %v", stream.StreamId, err)
			continue
		}
		calculator := metrics_service.NewCalculatorOfWaterLevels(stream.AlertLevel, stream.EvacuationLevel, stream.LowWaterLevel)
		for _, observedData := range values {
			if observedData.Value != nil {
				calculator.Compute(*observedData.Value)
				behaviourResponse.TotalValuesCount += 1
			}
		}
		behaviourResponse.CountAlertLevel += calculator.GetAlertsCount()
		behaviourResponse.CountLowWaterLevel += calculator.GetLowWaterCount()
		behaviourResponse.CountEvacuationLevel += calculator.GetEvacuationCount()
		levels := calculator.GetStreamLevels(stream.StreamId)
		if levels != nil && len(levels) > 0 {
			behaviourResponse.StreamLevels = append(behaviourResponse.StreamLevels, levels...)
		}
	}
	return behaviourResponse
}

func (s outputsServiceImpl) GetOutputBehaviourMetrics(configId uint64, timeStart time.Time, timeEnd time.Time) (*dtos.BehaviourStreamsResponse, error) {
	behaviourStreams, err := s.repository.GetStreamsForOutputMetrics(configId)
	if err != nil {
		return nil, err
	}

	return s.getLevelsCountForAllStreams(behaviourStreams, timeStart, timeEnd), nil
}
