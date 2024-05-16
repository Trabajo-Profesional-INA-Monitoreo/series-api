package jobs

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

func SetUpJobs(repositories *config.Repositories, apiConfig *config.ServiceConfigurationData) {
	log.Infof("Setting up jobs")
	c := cron.New()
	SetUpFaultDetectorJobs(repositories, apiConfig, c)
	SetUpNotificationsJobs(apiConfig, c)
	c.Start()
}
