package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const DayInHours = 24
const MinForecastedDays = 4

type FaultDetector interface {
	DetectFaults()
}

type StreamFaultDetector interface {
	handleStream(stream entities.Stream)
}

type faultDetectorService struct {
	streamsRepository            repositories.StreamRepository
	observedFaultDetectorService StreamFaultDetector
	forecastFaultDetectorService StreamFaultDetector
	detectionMaxThreads          int
}

func NewFaultDetectorService(streamsRepository repositories.StreamRepository,
	configuredStreamsRepository repositories.ConfiguredStreamsRepository,
	errorsRepository repositories.ErrorsRepository,
	inaApiClient clients.InaAPiClient,
	detectionMaxThreads int,
) FaultDetector {
	return &faultDetectorService{
		streamsRepository:            streamsRepository,
		detectionMaxThreads:          detectionMaxThreads,
		observedFaultDetectorService: newObservedFaultDetectorService(configuredStreamsRepository, errorsRepository, inaApiClient),
		forecastFaultDetectorService: newForecastFaultDetectorService(configuredStreamsRepository, errorsRepository, inaApiClient),
	}
}

func (f faultDetectorService) handleStream(streamChannel chan entities.Stream) error {
	stream := <-streamChannel
	if stream.IsObserved() {
		f.observedFaultDetectorService.handleStream(stream)
	} else if stream.IsForecasted() {
		f.forecastFaultDetectorService.handleStream(stream)
	}
	return nil
}

func (f faultDetectorService) DetectFaults() {
	log.Infof("FaultsDetector | Performing fault checks...")
	var goroutinePool errgroup.Group
	goroutinePool.SetLimit(f.detectionMaxThreads)
	streams := f.streamsRepository.GetStreams()
	for _, stream := range streams {
		streamChannel := make(chan entities.Stream, 1)
		streamChannel <- stream
		goroutinePool.Go(func() error { return f.handleStream(streamChannel) })
	}
	err := goroutinePool.Wait()
	if err != nil {
		log.Errorf("FaultsDetector | Error waiting for all goroutines to finish %v", err)
	}
	log.Infof("FaultsDetector | Fault check done")
}
