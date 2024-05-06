package jobs

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

func SetUpNotificationsJobs(apiConfig *config.ServiceConfigurationData, c *cron.Cron) {
	notificationsService := services.NewNotificationsService(
		clients.NewNotificationsAPiClientImpl(apiConfig),
	)

	err := c.AddFunc(apiConfig.DailyNotificationCronTime, notificationsService.SendDailyNotification)
	if err != nil {
		log.Fatalf("Error starting notifications service, stopping... | Err: %v", err)
	}

	log.Infof("Started notifications service with time: %v", apiConfig.DailyNotificationCronTime)
}
