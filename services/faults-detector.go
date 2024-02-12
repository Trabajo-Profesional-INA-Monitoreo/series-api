package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	log "github.com/sirupsen/logrus"
)

type FaultDetector interface {
	DetectFaults()
}

type faultDetectorService struct {
	streamsRepository repositories.StreamRepository
}

func NewFaultDetectorService(streamsRepository repositories.StreamRepository) FaultDetector {
	return &faultDetectorService{streamsRepository: streamsRepository}
}

func (f faultDetectorService) DetectFaults() {

	log.Debugf("Performing fault checks...")
	//streams := f.streamsRepository.GetStreams()
	//for _, stream := range streams {
	//
	//}
}
