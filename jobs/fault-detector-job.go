package jobs

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services/detection-services"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

func SetUpJobs(repositories *config.Repositories, apiConfig *config.ServiceConfigurationData) {
	c := cron.New()
	faultDetector := detection_services.NewFaultDetectorService(repositories.StreamsRepository,
		repositories.ConfiguredStreamRepository,
		repositories.ErrorsRepository,
		clients.NewInaApiClientImpl(apiConfig),
		apiConfig.DetectionMaxThreads,
	)
	entryId, err := c.AddFunc(apiConfig.FaultCronTime, faultDetector.DetectFaults)
	if err != nil {
		log.Fatalf("Error starting fault detector service, stopping... | Err: %v", err)
	}
	c.Start()
	log.Infof("Started fault detector service (Entry Id: %v) with time: %v", entryId, apiConfig.FaultCronTime)
}
