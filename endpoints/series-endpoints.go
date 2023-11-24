package endpoints

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/controllers"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func setSeriesEndpoints(apiGroup *gin.RouterGroup, config *config.ApiConfig) {
	controller := controllers.NewSeriesController(services.NewSeriesService(repositories.NewDbRepository(config.DbUrl)))
	testApi := apiGroup.Group("/series")
	{
		testApi.GET("/redes", controller.GetNetworks)
		testApi.GET("/estaciones", controller.GetStations)
	}
	log.Infof("Configured stream endpoints")
}
