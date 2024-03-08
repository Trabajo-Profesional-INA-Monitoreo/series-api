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
}

func NewFaultDetectorService(streamsRepository repositories.StreamRepository,
	configuredStreamsRepository repositories.ConfiguredStreamsRepository,
	errorsRepository repositories.ErrorsRepository,
	inaApiClient clients.InaAPiClient,
) FaultDetector {
	return &faultDetectorService{
		streamsRepository:           streamsRepository,
		configuredStreamsRepository: configuredStreamsRepository,
		errorsRepository:            errorsRepository,
		inaApiClient:                inaApiClient,
	}
}

func (f faultDetectorService) handleMissingForecast(stream entities.Stream, configuredStream entities.ConfiguredStream, res *responses.LastForecast) {
	now := time.Now()
	diff := now.Sub(res.ForecastDate)
	reqErrorId := fmt.Sprintf("%v", res.RunId)
	detectedError := f.errorsRepository.GetDetectedErrorForStreamWithIdAndType(stream.StreamId, reqErrorId, entities.ForecastMissing)
	detected := detectedError.RequestId == reqErrorId
	if diff.Hours() > configuredStream.UpdateFrequency && !detected {
		// There should be a new forecast already
		// We save the detected error
		log.Debugf("Detected missing forecast for: %v", stream.StreamId)
		detected := entities.DetectedError{
			StreamId:         stream.StreamId,
			Stream:           &stream,
			ConfiguredStream: []entities.ConfiguredStream{configuredStream},
			DetectedDate:     time.Now(),
			RequestId:        reqErrorId,
			ErrorType:        entities.ForecastMissing,
		}
		f.errorsRepository.Create(detected)
	} else if detected {
		// We already detected the error, we need to save the relationship to the current ConfiguredStream
		detectedError := f.errorsRepository.GetDetectedErrorForStreamWithIdAndType(stream.StreamId, reqErrorId, entities.ForecastMissing)
		detectedError.ConfiguredStream = append(detectedError.ConfiguredStream, configuredStream)
		f.errorsRepository.Update(detectedError)
	}
}

func (f faultDetectorService) handleNullValues(data []responses.ObservedDataResponse, stream entities.Stream, configuredStreams []entities.ConfiguredStream) {
	for _, observed := range data {
		isNull := observed.Value == 0 || observed.DataId == ""
		reqId := fmt.Sprintf("%v", observed.TimeStart)
		if isNull && !f.errorsRepository.AlreadyDetectedErrorForStreamWithIdAndType(stream.StreamId, reqId, entities.NullValue) {
			// We detected a new null value
			log.Debugf("Detected null value for: %v", stream.StreamId)
			detected := entities.DetectedError{
				StreamId:         stream.StreamId,
				Stream:           &stream,
				ConfiguredStream: configuredStreams,
				DetectedDate:     time.Now(),
				RequestId:        reqId,
				ErrorType:        entities.NullValue,
			}
			f.errorsRepository.Create(detected)
		}
	}
}

func (f faultDetectorService) getObservedDataFromStream(streamId uint64) ([]responses.ObservedDataResponse, error) {
	res, err := f.inaApiClient.GetObservedData(streamId, time.Now().Add(time.Duration(-24)*time.Hour), time.Now())
	if err != nil {
		return nil, fmt.Errorf("error performing check for stream %v: %v", streamId, err)
	}
	return res, nil
}

func (f faultDetectorService) handleObservedStreams(stream entities.Stream) {
	configuredStreams := f.configuredStreamsRepository.FindConfiguredStreamsWithCheckErrorsForStream(stream)
	data, err := f.getObservedDataFromStream(stream.StreamId)
	if err != nil {
		log.Errorf("Error detecting observed errors: %v", err)
		return
	}
	f.handleNullValues(data, stream, configuredStreams)
}

func (f faultDetectorService) handleForecastedStream(stream entities.Stream) {
	configuredStreams := f.configuredStreamsRepository.FindConfiguredStreamsWithCheckErrorsForStream(stream)
	for _, configuredStream := range configuredStreams {
		res, err := f.inaApiClient.GetLastForecast(configuredStream.CalibrationId)
		if err != nil {
			log.Errorf("Error performing check for stream %v with configuration %v", stream.StreamId, configuredStream.ConfiguredStreamId)
			return
		}
		f.handleMissingForecast(stream, configuredStream, res)

		// TODO handle 4 days forecast horizon
	}
}

func (f faultDetectorService) DetectFaults() {

	log.Debugf("Performing fault checks...")
	streams := f.streamsRepository.GetStreams()
	for _, stream := range streams {
		if stream.IsObserved() {
			f.handleObservedStreams(stream)
		} else if stream.IsForecasted() {
			f.handleForecastedStream(stream)
		}

	}
}
