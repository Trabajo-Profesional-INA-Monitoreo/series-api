package services

import (
	"errors"
	"strconv"
	"time"

	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/converters"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	exceptions "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/errors"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	log "github.com/sirupsen/logrus"
)

type StreamService interface {
	GetStations(time.Time, time.Time, uint64) dtos.StreamsPerStationResponse
	GetCuredSerieById(id string, start time.Time, end time.Time) dtos.StreamsDataResponse
	GetObservatedSerieById(id string, start time.Time, end time.Time) dtos.StreamsDataResponse
	GetPredictedSerieById(id string, streamId uint64) dtos.CalibratedStreamsDataResponse
	GetStreamData(streamId uint64, configId uint64, timeStart time.Time, timeEnd time.Time) (*dtos.StreamData, error)
	CreateStream(streamId uint64, streamType uint64) error
	GetStreamCards(parameters *dtos.QueryParameters) (*dtos.StreamCardsResponse, error)
	GetOutputBehaviourMetrics(configId uint64, timeStart time.Time, timeEnd time.Time) (*dtos.BehaviourStreamsResponse, error)
	GetNodes(start time.Time, end time.Time, id uint64) dtos.StreamsPerNodeResponse
	GetRedundancies(configuredStreamId string) dtos.Redundancies
}

type streamService struct {
	repository                  repositories.StreamRepository
	inaApiClient                clients.InaAPiClient
	configuredStreamsRepository repositories.ConfiguredStreamsRepository
	nodesRepository             repositories.NodeRepository
}

func NewStreamService(repository repositories.StreamRepository, inaApiClient clients.InaAPiClient, configuredStreamsRepository repositories.ConfiguredStreamsRepository, nodeRepository repositories.NodeRepository) StreamService {
	return &streamService{repository, inaApiClient, configuredStreamsRepository, nodeRepository}
}

func (s streamService) GetNodes(timeStart time.Time, timeEnd time.Time, configId uint64) dtos.StreamsPerNodeResponse {
	nodes := s.nodesRepository.GetStreamsPerNodeById(strconv.FormatUint(configId, 10))
	errorsPerNode := s.repository.GetErrorsOfNodes(configId, timeStart, timeEnd)

	for _, errors := range errorsPerNode {
		for _, node := range nodes {
			if node.NodeId == errors.NodeId {
				node.ErrorCount = errors.ErrorCount
				break
			}
		}
	}
	return dtos.StreamsPerNodeResponse{Nodes: nodes}
}

func (s streamService) GetStations(timeStart time.Time, timeEnd time.Time, configId uint64) dtos.StreamsPerStationResponse {
	stations := s.repository.GetStations(configId)
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

func (s streamService) GetCuredSerieById(id string, start time.Time, end time.Time) dtos.StreamsDataResponse {
	num, _ := strconv.ParseUint(id, 10, 64)
	streams, _ := s.inaApiClient.GetObservedData(num, start, end)
	var streamsData []dtos.StreamsData
	for _, stream := range streams {
		streamsData = append(streamsData, stream.ConvertToStreamData())
	}
	return dtos.StreamsDataResponse{Streams: streamsData}
}

func (s streamService) GetObservatedSerieById(id string, start time.Time, end time.Time) dtos.StreamsDataResponse {
	num, _ := strconv.ParseUint(id, 10, 64)
	streams, _ := s.inaApiClient.GetObservedData(num, start, end)
	var streamsData []dtos.StreamsData
	for _, stream := range streams {
		streamsData = append(streamsData, stream.ConvertToStreamData())
	}
	return dtos.StreamsDataResponse{Streams: streamsData}
}

func (s streamService) GetPredictedSerieById(id string, streamId uint64) dtos.CalibratedStreamsDataResponse {
	num, _ := strconv.ParseUint(id, 10, 64)
	streams, _ := s.inaApiClient.GetLastForecast(num)
	return streams.ConvertToCalibratedStreamsDataResponse(streamId)
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
		return getMetricsForForecastedStream(values, neededMetrics, waterLevelCalculator), &values.ForecastDate
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

func (s streamService) CreateStream(streamId uint64, streamType uint64) error {
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

func (s streamService) GetOutputBehaviourMetrics(configId uint64, timeStart time.Time, timeEnd time.Time) (*dtos.BehaviourStreamsResponse, error) {
	behaviourStreams, err := s.repository.GetStreamsForOutputMetrics(configId)
	if err != nil {
		return nil, err
	}

	return getLevelsCountForAllStreams(behaviourStreams, timeStart, timeEnd, s.inaApiClient), nil
}

func (s streamService) GetRedundancies(configuredStreamId string) dtos.Redundancies {
	redundancies := s.repository.GetRedundancies(configuredStreamId)

	return dtos.Redundancies{Redundancies: redundancies}
}
