package endpoints

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/controllers"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func setErrorEndpoints(apiGroup *gin.RouterGroup, errorsRepository repositories.ErrorsRepository) {
	controller := controllers.NewErrorsController(services.NewErrorsService(errorsRepository))
	testApi := apiGroup.Group("/errores")
	{
		testApi.GET("/por-dia", controller.GetErrorsPerDay)
	}
	log.Infof("Configured error endpoints")
}