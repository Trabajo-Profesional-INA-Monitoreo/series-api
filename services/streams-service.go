package services

import (
	"errors"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"time"

	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/converters"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	exceptions "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/errors"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	log "github.com/sirupsen/logrus"
)

type StreamCreationService interface {
	CreateStream(streamId uint64, streamType entities.StreamType) error
}

type StreamRetrievalService interface {
	GetStreamData(streamId uint64, configId uint64, timeStart time.Time, timeEnd time.Time) (*dtos.StreamData, error)
	GetStreamCards(parameters *dtos.QueryParameters) (*dtos.StreamCardsResponse, error)
	GetRedundancies(configuredStreamId uint64) dtos.Redundancies
}

type StreamService interface {
	StreamCreationService
	StreamRetrievalService
}

type streamService struct {
	repository                  repositories.ManagerStreamsRepository
	inaApiClient                clients.InaAPiClient
	configuredStreamsRepository repositories.MetricsConfiguredStreamsRepository
}

func NewStreamService(repository *config.Repositories, inaApiClient clients.InaAPiClient) StreamService {
	return &streamService{repository.StreamsRepository, inaApiClient, repository.ConfiguredStreamRepository}
}

func (s streamService) getMetricsFromConfiguredStream(stream entities.Stream, configured entities.ConfiguredStream, timeStart time.Time, timeEnd time.Time) (*[]dtos.MetricCard, *time.Time) {
	neededMetrics := configured.Metrics
	waterLevelCalculator := NewCalculatorOfWaterLevelsDependingOnVariable(*stream.Station, stream.VariableId)
	if stream.IsForecasted() {
		values, err := s.inaApiClient.GetLastForecast(configured.CalibrationId)
		if err != nil {
			log.Errorf("Could not get metrics with calibration id %v: %v", configured.CalibrationId, err)
			return nil, nil
		}
		forecast := values.GetForecastOfStream(stream.StreamId)
		return getMetricsForForecastedStream(forecast, neededMetrics, waterLevelCalculator), &values.ForecastDate
	}
	lastUpdateResponse := make(chan *time.Time)
	go func() {
		data, err := s.inaApiClient.GetStream(stream.StreamId)
		defer close(lastUpdateResponse)
		if err != nil {
			lastUpdateResponse <- nil
			return
		}
		lastUpdateResponse <- data.DateRange.TimeEnd
	}()
	values, err := s.inaApiClient.GetObservedData(stream.StreamId, timeStart, timeEnd)
	if err != nil {
		log.Errorf("Could not get metrics for stream with id %v: %v", stream.StreamId, err)
		return nil, nil
	}

	return getMetricsForObservedOrCuratedStream(values, neededMetrics, waterLevelCalculator), <-lastUpdateResponse
}

func (s streamService) GetStreamData(streamId uint64, configId uint64, timeStart time.Time, timeEnd time.Time) (*dtos.StreamData, error) {
	stream, err := s.repository.GetStreamWithAssociatedData(streamId)
	if err != nil {
		return nil, err
	}
	configured, err := s.configuredStreamsRepository.FindConfiguredStreamById(configId)
	if err != nil {
		return nil, err
	}
	streamData := dtos.NewStreamData(stream, configured)
	metrics, lastUpdate := s.getMetricsFromConfiguredStream(stream, configured, timeStart, timeEnd)
	streamData.Metrics = metrics
	streamData.LastUpdate = lastUpdate
	return streamData, nil
}

func (s streamService) CreateStream(streamId uint64, streamType entities.StreamType) error {
	_, err := s.repository.GetStreamWithAssociatedData(streamId)
	if errors.Is(err, &exceptions.NotFound{}) {

		inaStreamResponse, err := s.inaApiClient.GetStream(streamId)
		if err != nil {
			return err
		}

		err = s.repository.CreateUnit(converters.ConvertUnitResponseToEntity(inaStreamResponse.Unit))
		if err != nil {
			return err
		}

		err = s.repository.CreateStation(converters.ConvertStationResponseToEntity(inaStreamResponse.Station))
		if err != nil {
			return err
		}

		err = s.repository.CreateProcedure(converters.ConvertProcedureResponseToEntity(inaStreamResponse.Procedure))
		if err != nil {
			return err
		}

		err = s.repository.CreateVariable(converters.ConvertVariableResponseToEntity(inaStreamResponse.Variable))
		if err != nil {
			return err
		}

		err = s.repository.CreateStream(converters.ConvertStreamResponseToEntity(inaStreamResponse, streamType))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s streamService) GetStreamCards(parameters *dtos.QueryParameters) (*dtos.StreamCardsResponse, error) {
	result, err := s.repository.GetStreamCards(*parameters)
	if err != nil {
		return nil, err
	}
	var configuredIds []uint64
	for _, card := range *result.Content {
		if card.CheckErrors {
			configuredIds = append(configuredIds, card.ConfiguredStreamId)
		}
	}
	errorsPerConfigStream, err := s.configuredStreamsRepository.CountErrorOfConfigurations(configuredIds, parameters)
	if err != nil {
		return nil, err
	}
	for _, errors := range errorsPerConfigStream {
		for _, card := range *result.Content {
			if errors.ConfiguredStreamId == card.ConfiguredStreamId {
				card.TotalErrors = &errors.ErrorsCount
				break
			}
		}
	}
	return result, nil
}

func (s streamService) GetRedundancies(configuredStreamId uint64) dtos.Redundancies {
	redundancies := s.repository.GetRedundancies(configuredStreamId)

	return dtos.Redundancies{Redundancies: redundancies}
}
