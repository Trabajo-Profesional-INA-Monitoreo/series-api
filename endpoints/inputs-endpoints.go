package endpoints

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/controllers"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func setInputsEndpoints(apiGroup *gin.RouterGroup, config *config.ApiConfig) {
	controller := controllers.NewInputsController(services.NewInputsService(repositories.NewDbRepository(config.DbUrl)))
	testApi := apiGroup.Group("/inputs")
	{
		testApi.GET("/metricas-generales", controller.GetGeneralMetrics)
	}
	log.Infof("Configured inputs endpoints")
}
