package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type StreamService interface {
	GetNetworks() dtos.StreamsPerNetworkResponse
	GetStations(time.Time, time.Time, uint64) dtos.StreamsPerStationResponse
	GetCuredSerieById(id string, start time.Time, end time.Time) dtos.StreamsDataResponse
	GetObservatedSerieById(id string, start time.Time, end time.Time) dtos.StreamsDataResponse
	GetPredictedSerieById(id string) dtos.CalibratedStreamsDataResponse
	GetStreamData(streamId uint64, configId uint64, timeStart time.Time, timeEnd time.Time) (*dtos.StreamData, error)
	GetStreamCards(parameters *dtos.StreamCardsParameters) (*dtos.StreamCardsResponse, error)
	GetOutputBehaviourMetrics(configId uint64, timeStart time.Time, timeEnd time.Time) (*dtos.BehaviourStreamsResponse, error)
}

type streamService struct {
	repository                  repositories.StreamRepository
	inaApiClient                clients.InaAPiClient
	configuredStreamsRepository repositories.ConfiguredStreamsRepository
}

func NewSeriesService(repository repositories.StreamRepository, inaApiClient clients.InaAPiClient, configuredStreamsRepository repositories.ConfiguredStreamsRepository) StreamService {
	return &streamService{repository, inaApiClient, configuredStreamsRepository}
}

func (s streamService) GetNetworks() dtos.StreamsPerNetworkResponse {
	networks := s.repository.GetNetworks()
	return dtos.StreamsPerNetworkResponse{Networks: networks}
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

func (s streamService) GetPredictedSerieById(id string) dtos.CalibratedStreamsDataResponse {
	num, _ := strconv.ParseUint(id, 10, 64)
	streams, _ := s.inaApiClient.GetLastForecast(num)
	return streams.ConvertToCalibratedStreamsDataResponse()
}

func (s streamService) getMetricsFromConfiguredStream(stream entities.Stream, configured entities.ConfiguredStream, timeStart time.Time, timeEnd time.Time) *[]dtos.MetricCard {
	neededMetrics := configured.Metrics
	if len(neededMetrics) == 0 {
		return nil
	}
	waterLevelCalculator := NewCalculatorOfWaterLevelsDependingOnVariable(*stream.Station, stream.VariableId)
	if stream.IsForecasted() {
		values, err := s.inaApiClient.GetLastForecast(configured.CalibrationId)
		if err != nil {
			log.Errorf("Could not get metrics with calibration id %v: %v", configured.CalibrationId, err)
			return nil
		}
		return getMetricsForForecastedStream(values, neededMetrics, waterLevelCalculator)
	}
	values, err := s.inaApiClient.GetObservedData(stream.StreamId, timeStart, timeEnd)
	if err != nil {
		log.Errorf("Could not get metrics for stream with id %v: %v", stream.StreamId, err)
		return nil
	}
	return getMetricsForObservedOrCuratedStream(values, neededMetrics, waterLevelCalculator)
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
	streamData.Metrics = s.getMetricsFromConfiguredStream(stream, configured, timeStart, timeEnd)
	return streamData, nil
}

func (s streamService) GetStreamCards(parameters *dtos.StreamCardsParameters) (*dtos.StreamCardsResponse, error) {
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
