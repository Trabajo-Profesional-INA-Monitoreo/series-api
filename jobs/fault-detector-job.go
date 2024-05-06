package jobs

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services/detection-services"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

func SetUpFaultDetectorJobs(repositories *config.Repositories, apiConfig *config.ServiceConfigurationData, c *cron.Cron) {

	faultDetector := detection_services.NewFaultDetectorService(repositories.StreamsRepository,
		repositories.ConfiguredStreamRepository,
		repositories.ErrorsRepository,
		clients.NewInaApiClientImpl(apiConfig),
		apiConfig.DetectionMaxThreads,
	)
	err := c.AddFunc(apiConfig.FaultCronTime, faultDetector.DetectFaults)
	if err != nil {
		log.Fatalf("Error starting fault detector service, stopping... | Err: %v", err)
	}
	log.Infof("Started fault detector service with time: %v", apiConfig.FaultCronTime)
}
