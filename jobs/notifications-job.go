package jobs

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"time"
)

func SetUpNotificationsJobs(apiConfig *config.ServiceConfigurationData) {
	customLocation, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		log.Fatalf("Error starting notifications service, stopping... | Err: %v", err)
	}
	c := cron.NewWithLocation(customLocation)

	notificationsService := services.NewNotificationsService(
		clients.NewNotificationsAPiClientImpl(apiConfig),
	)

	err = c.AddFunc(apiConfig.DailyNotificationCronTime, notificationsService.SendDailyNotification)
	if err != nil {
		log.Fatalf("Error starting notifications service, stopping... | Err: %v", err)
	}

	c.Start()

	log.Infof("Started notifications service with time: %v", apiConfig.DailyNotificationCronTime)
}
