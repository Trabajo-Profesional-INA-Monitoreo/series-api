package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
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
	inaApiClient                clients.InaAPiClient
}

func NewFaultDetectorService(streamsRepository repositories.StreamRepository,
	configuredStreamsRepository repositories.ConfiguredStreamsRepository,
	inaApiClient clients.InaAPiClient) FaultDetector {
	return &faultDetectorService{
		streamsRepository:           streamsRepository,
		configuredStreamsRepository: configuredStreamsRepository,
		inaApiClient:                inaApiClient,
	}
}

func (f faultDetectorService) handleForecastedStream(stream entities.Stream, configuredStream entities.ConfiguredStream) {
	res, err := f.inaApiClient.GetLastForecast(configuredStream.CalibrationId)
	if err != nil {
		log.Errorf("Error performing check for stream %v with configuration %v", stream.StreamId, configuredStream.ConfiguredStreamId)
		return
	}
	now := time.Now()
	diff := now.Sub(res.ForecastDate)
	if diff.Hours() > 6 {
		// There should be a new forecast already
	}
}

func (f faultDetectorService) DetectFaults() {

	log.Debugf("Performing fault checks...")
	streams := f.streamsRepository.GetStreams()
	for _, stream := range streams {
		configuredStreams := f.configuredStreamsRepository.FindConfiguredStreamsForStream(stream)
		for _, configuredStream := range configuredStreams {

			if stream.IsForecasted() {
				f.handleForecastedStream(stream, configuredStream)
			}
		}
	}
}
