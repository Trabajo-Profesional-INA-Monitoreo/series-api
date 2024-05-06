package jobs

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"time"
)

func SetUpJobs(repositories *config.Repositories, apiConfig *config.ServiceConfigurationData) {
	log.Infof("Setting up jobs")
	customLocation, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		log.Fatalf("Error starting notifications service, stopping... | Err: %v", err)
	}
	c := cron.NewWithLocation(customLocation)
	SetUpFaultDetectorJobs(repositories, apiConfig, c)
	SetUpNotificationsJobs(apiConfig, c)
	c.Start()
}
