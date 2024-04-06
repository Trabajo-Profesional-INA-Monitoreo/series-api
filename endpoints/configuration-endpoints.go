package endpoints

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/controllers"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func setConfigurationEndpoints(apiGroup *gin.RouterGroup, repositories *config.Repositories, inaClient clients.InaAPiClient) {
	controller := controllers.NewConfigurationController(services.NewConfigurationService(repositories, inaClient))
	testApi := apiGroup.Group("/configuracion")
	{
		testApi.GET("", controller.GetAllConfigurations)
		testApi.GET("/:id", controller.GetConfigurationById)
		testApi.POST("", controller.CreateConfiguration)
		testApi.PUT("", controller.ModifyConfiguration)
		testApi.DELETE("/:id", controller.DeleteConfiguration)
	}
	log.Infof("Configured configuration endpoints")
}
