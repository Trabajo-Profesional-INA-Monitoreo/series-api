package endpoints

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/controllers"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func setConfigurationEndpoints(apiGroup *gin.RouterGroup, configurationRepository repositories.ConfigurationRepository) {
	controller := controllers.NewConfigurationController(services.NewConfigurationService(configurationRepository))
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
