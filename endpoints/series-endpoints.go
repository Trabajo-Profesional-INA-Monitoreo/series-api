package endpoints

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/controllers"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func setSeriesEndpoints(apiGroup *gin.RouterGroup, streamsRepository repositories.StreamRepository, inaApiClient clients.InaAPiClient) {
	controller := controllers.NewSeriesController(services.NewSeriesService(streamsRepository, inaApiClient))
	testApi := apiGroup.Group("/series")
	{
		testApi.GET("/redes", controller.GetNetworks)
		testApi.GET("/estaciones", controller.GetStations)
		testApi.GET("/curadas/:serie_id", controller.GetCuredSerieById)
		testApi.GET("/observadas/:serie_id", controller.GetObservatedSerieById)
		testApi.GET("/pronosticadas/:calibrado_id", controller.GetPredictedSerieById)
	}
	log.Infof("Configured stream endpoints")
}
