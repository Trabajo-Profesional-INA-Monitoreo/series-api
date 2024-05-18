package endpoints

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/controllers"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services/inputs_service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func setInputsEndpoints(apiGroup *gin.RouterGroup, inputsRepository repositories.InputsRepository) {
	controller := controllers.NewInputsController(inputs_service.NewInputsService(inputsRepository))
	testApi := apiGroup.Group("/inputs")
	{
		testApi.GET("/metricas-generales", controller.GetGeneralMetrics)
		testApi.GET("/series-con-nulos", controller.GetTotalStreamsWithNullValues)
		testApi.GET("/series-fuera-umbral", controller.GetTotalStreamsWithObservedOutlier)
		testApi.GET("/series-retardos", controller.GetTotalStreamsWithDelay)
	}
	log.Infof("Configured inputs endpoints")
}
