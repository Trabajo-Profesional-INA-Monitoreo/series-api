package endpoints

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/controllers"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func setInputsEndpoints(apiGroup *gin.RouterGroup, streamsRepository repositories.StreamRepository) {
	controller := controllers.NewInputsController(services.NewInputsService(streamsRepository))
	testApi := apiGroup.Group("/inputs")
	{
		testApi.GET("/metricas-generales", controller.GetGeneralMetrics)
		testApi.GET("/series-con-nulos", controller.GetTotalStreamsWithNullValues)
		testApi.GET("/series-fuera-umbral", controller.GetTotalStreamsWithObservedOutlier)
	}
	log.Infof("Configured inputs endpoints")
}
