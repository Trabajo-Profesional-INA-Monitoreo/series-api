package endpoints

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/controllers"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func SetNotificationEndpoints(apiGroup *gin.RouterGroup, config *config.ServiceConfigurationData) {
	controller := controllers.NewNotificationController(clients.NewNotificationsAPiClientImpl(config))
	testApi := apiGroup.Group("/notificacion")
	{
		testApi.POST("", controller.CreateNotification)
	}
	log.Infof("Configured notification endpoints")
}
