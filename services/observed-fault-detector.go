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

type observedFaultDetectorService struct {
	configuredStreamsRepository repositories.DetectionConfiguredStreamsRepository
	errorsRepository            repositories.ErrorDetectionRepository
	inaApiClient                clients.InaAPiClient
}

func newObservedFaultDetectorService(configuredStreamsRepository repositories.DetectionConfiguredStreamsRepository, errorsRepository repositories.ErrorDetectionRepository, inaApiClient clients.InaAPiClient) StreamFaultDetector {
	return &observedFaultDetectorService{configuredStreamsRepository: configuredStreamsRepository, errorsRepository: errorsRepository, inaApiClient: inaApiClient}
}

func (f observedFaultDetectorService) checkOutliers(data []responses.ObservedDataResponse, stream entities.Stream, configuredStream entities.ConfiguredStream) {
	consecutiveOutliers := 0
	for _, observed := range data {
		isOutlier := valueIsAnOutlier(configuredStream, observed)
		if isOutlier && consecutiveOutliers == 0 {
			reqId := fmt.Sprintf("%v", observed.TimeStart)
			detectedError := f.errorsRepository.GetDetectedErrorForStreamWithIdAndType(stream.StreamId, reqId, entities.ObservedOutlier)
			detected := detectedError.RequestId == reqId
			if !detected {
				log.Debugf("Detected outlier value in observed stream for: %v", stream.StreamId)
				detected := entities.DetectedError{
					StreamId:         stream.StreamId,
					Stream:           &stream,
					ConfiguredStream: []entities.ConfiguredStream{configuredStream},
					DetectedDate:     time.Now(),
					RequestId:        reqId,
					ErrorType:        entities.ObservedOutlier,
					ExtraInfo:        fmt.Sprintf("Valor %v - Timestamp %v", *observed.Value, observed.TimeStart),
				}
				f.errorsRepository.Create(detected)
			} else if !contains(detectedError.ConfiguredStream, configuredStream) {
				detectedError.ConfiguredStream = append(detectedError.ConfiguredStream, configuredStream)
				f.errorsRepository.Update(detectedError)
			}
			consecutiveOutliers++
		}
		if !isOutlier {
			consecutiveOutliers = 0
		}
	}
}

func (f observedFaultDetectorService) checkNullValues(data []responses.ObservedDataResponse, stream entities.Stream, configuredStreams []entities.ConfiguredStream) {
	for _, observed := range data {
		isNull := observed.Value == nil || observed.DataId == ""
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
				ExtraInfo:        fmt.Sprintf("Timestamp del dato nulo %v", observed.TimeStart),
			}
			f.errorsRepository.Create(detected)
		}
	}
}

func (f observedFaultDetectorService) handleStream(stream entities.Stream) {
	configuredStreams := f.configuredStreamsRepository.FindConfiguredStreamsWithCheckErrorsForStream(stream)
	data, err := getObservedDataFromStream(stream.StreamId, time.Now().Add(time.Duration(-24)*time.Hour), time.Now(), f.inaApiClient)
	if err != nil {
		log.Errorf("Error detecting observed errors: %v", err)
		return
	}
	f.checkNullValues(data, stream, configuredStreams)
	for _, configuredStream := range configuredStreams {
		f.checkOutliers(data, stream, configuredStream)
	}
}
