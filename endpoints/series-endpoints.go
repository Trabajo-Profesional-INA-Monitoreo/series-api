package endpoints

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/controllers"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services/nodes-service"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services/outputs-service"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services/stations-service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func setSeriesEndpoints(apiGroup *gin.RouterGroup, repositories *config.Repositories, inaApiClient clients.InaAPiClient) {
	streamsService := services.NewStreamService(repositories, inaApiClient)
	streamsController := controllers.NewSeriesController(streamsService)
	outputsController := controllers.NewOutputsController(outputs_service.NewOutputsService(repositories.OutputsRepository, inaApiClient))
	inaController := controllers.NewInaController(services.NewInaServiceApi(inaApiClient))
	nodesController := controllers.NewNodesController(nodes_service.NewNodesService(repositories, inaApiClient))
	stationsController := controllers.NewStationsController(stations_service.NewStationsService(repositories.StationsRepository, inaApiClient))
	testApi := apiGroup.Group("/series")
	{
		testApi.GET("", streamsController.GetStreamCards)
		testApi.GET("/:serie_id", streamsController.GetStreamDataById)
		testApi.GET("/redundancias/:configured_stream_id", streamsController.GetRedundancies)
		testApi.GET("/nodos", nodesController.GetNodes)
		testApi.GET("/estaciones", stationsController.GetStations)
		testApi.GET("/comportamiento", outputsController.GetOutputMetrics)
		testApi.GET("/observadas/:serie_id", inaController.GetObservatedSerieById)
		testApi.GET("/pronosticadas/:calibrado_id", inaController.GetPredictedSerieById)
	}
	log.Infof("Configured stream endpoints")
}
