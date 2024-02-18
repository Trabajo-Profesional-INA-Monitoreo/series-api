package services

import (
	"fmt"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients/responses"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	log "github.com/sirupsen/logrus"
	"time"
)

type FaultDetector interface {
	DetectFaults()
}

type faultDetectorService struct {
	streamsRepository           repositories.StreamRepository
	configuredStreamsRepository repositories.ConfiguredStreamsRepository
	errorsRepository            repositories.ErrorsRepository
	inaApiClient                clients.InaAPiClient
	forecastMaxWaitingTimeHours float64
}

func NewFaultDetectorService(streamsRepository repositories.StreamRepository,
	configuredStreamsRepository repositories.ConfiguredStreamsRepository,
	errorsRepository repositories.ErrorsRepository,
	inaApiClient clients.InaAPiClient,
	forecastMaxWaitingTimeHours float64,
) FaultDetector {
	return &faultDetectorService{
		streamsRepository:           streamsRepository,
		configuredStreamsRepository: configuredStreamsRepository,
		errorsRepository:            errorsRepository,
		inaApiClient:                inaApiClient,
		forecastMaxWaitingTimeHours: forecastMaxWaitingTimeHours,
	}
}

func (f faultDetectorService) handleForecastedStream(stream entities.Stream, configuredStream entities.ConfiguredStream) {
	res, err := f.inaApiClient.GetLastForecast(configuredStream.CalibrationId)
	if err != nil {
		log.Errorf("Error performing check for stream %v with configuration %v", stream.StreamId, configuredStream.ConfiguredStreamId)
		return
	}
	f.handleMissingForecast(stream, configuredStream, res)
	// TODO handle 4 days forecast horizon
}

func (f faultDetectorService) handleMissingForecast(stream entities.Stream, configuredStream entities.ConfiguredStream, res *responses.LastForecast) {
	now := time.Now()
	diff := now.Sub(res.ForecastDate)
	if diff.Hours() > f.forecastMaxWaitingTimeHours && !f.errorsRepository.AlreadyDetectedErrorWithIdAndType(fmt.Sprintf("%v", res.RunId), entities.ForecastMissing) {
		// There should be a new forecast already
		// We save the detected error
		log.Debugf("Detected missing forecast for: %v", stream.StreamId)
		detected := entities.DetectedError{
			StreamId:           stream.StreamId,
			Stream:             &stream,
			ConfiguredStreamId: configuredStream.ConfiguredStreamId,
			ConfiguredStream:   &configuredStream,
			DetectedDate:       time.Now(),
			RequestId:          fmt.Sprintf("%v", res.RunId),
			ErrorType:          entities.ForecastMissing,
		}
		f.errorsRepository.Save(detected)
	}
}

func (f faultDetectorService) handleNullValues(stream entities.Stream, configuredStream entities.ConfiguredStream) {
	res, err := f.inaApiClient.GetObservedData(stream.StreamId, time.Now().Add(time.Duration(-24)*time.Hour), time.Now())
	if err != nil {
		log.Errorf("Error performing check for stream %v with configuration %v", stream.StreamId, configuredStream.ConfiguredStreamId)
		return
	}
	for _, observed := range res {
		isNull := observed.Value == 0 || observed.DataId == ""
		if isNull && !f.errorsRepository.AlreadyDetectedErrorWithIdAndType(fmt.Sprintf("%v", observed.TimeStart), entities.NullValue) {
			// We detected a new null value
			log.Debugf("Detected null value for: %v", stream.StreamId)
			detected := entities.DetectedError{
				StreamId:           stream.StreamId,
				Stream:             &stream,
				ConfiguredStreamId: configuredStream.ConfiguredStreamId,
				ConfiguredStream:   &configuredStream,
				DetectedDate:       time.Now(),
				RequestId:          fmt.Sprintf("%v", observed.TimeStart),
				ErrorType:          entities.NullValue,
			}
			f.errorsRepository.Save(detected)
		}
	}
}

func (f faultDetectorService) DetectFaults() {

	log.Debugf("Performing fault checks...")
	streams := f.streamsRepository.GetStreams()
	for _, stream := range streams {
		configuredStreams := f.configuredStreamsRepository.FindConfiguredStreamsForStream(stream)
		for _, configuredStream := range configuredStreams {
			if stream.IsObserved() {
				f.handleNullValues(stream, configuredStream)
			}
			if stream.IsForecasted() {
				f.handleForecastedStream(stream, configuredStream)
			}
		}
	}
}
