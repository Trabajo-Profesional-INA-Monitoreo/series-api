package jobs

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	log "github.com/sirupsen/logrus"
)

func SetUpJobs(repositories *config.Repositories, apiConfig *config.ServiceConfigurationData) {
	log.Infof("Setting up jobs")
	SetUpFaultDetectorJobs(repositories, apiConfig)
	SetUpNotificationsJobs(apiConfig)
}
