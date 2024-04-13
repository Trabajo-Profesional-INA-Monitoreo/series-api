package endpoints

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/controllers"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func setSeriesEndpoints(apiGroup *gin.RouterGroup, repositories *config.Repositories, inaApiClient clients.InaAPiClient) {
	controller := controllers.NewSeriesController(services.NewStreamService(repositories.StreamsRepository, inaApiClient, repositories.ConfiguredStreamRepository))
	testApi := apiGroup.Group("/series")
	{
		testApi.GET("", controller.GetStreamCards)
		testApi.GET("/comportamiento", controller.GetOutputMetrics)
		testApi.GET("/estaciones", controller.GetStations)
		testApi.GET("/:serie_id", controller.GetStreamDataById)
		testApi.GET("/curadas/:serie_id", controller.GetCuredSerieById)
		testApi.GET("/observadas/:serie_id", controller.GetObservatedSerieById)
		testApi.GET("/pronosticadas/:calibrado_id", controller.GetPredictedSerieById)
	}
	log.Infof("Configured stream endpoints")
}
