package endpoints

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/controllers"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/middlewares"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func setConfigurationEndpoints(apiGroup *gin.RouterGroup, repositories *config.Repositories, inaClient clients.InaAPiClient, apiConfig *config.ApiConfig) {
	controller := controllers.NewConfigurationController(services.NewConfigurationService(repositories, inaClient))
	testApi := apiGroup.Group("/configuracion")
	{
		testApi.GET("", controller.GetAllConfigurations)
		testApi.GET("/:id", controller.GetConfigurationById)
		testApi.POST("", middlewares.IsAnAdminToken(apiConfig), controller.CreateConfiguration)
		testApi.PUT("", middlewares.IsAnAdminToken(apiConfig), controller.ModifyConfiguration)
		testApi.DELETE("/:id", middlewares.IsAnAdminToken(apiConfig), controller.DeleteConfiguration)
	}
	log.Infof("Configured configuration endpoints")
}
