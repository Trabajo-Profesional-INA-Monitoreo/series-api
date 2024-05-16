package endpoints

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/controllers"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services/metrics-service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func setErrorEndpoints(apiGroup *gin.RouterGroup, errorsRepository repositories.ErrorsRepository) {
	controller := controllers.NewErrorsController(metrics_service.NewErrorsService(errorsRepository))
	testApi := apiGroup.Group("/errores")
	{
		testApi.GET("/por-dia", controller.GetErrorsPerDay)
		testApi.GET("/indicadores", controller.GetErrorIndicators)
		testApi.GET("/series-implicadas", controller.GetStreamsWithRelatedError)
		testApi.GET("/:configuredStreamId", controller.GetErrorsOfConfiguredStream)
	}
	log.Infof("Configured error endpoints")
}
