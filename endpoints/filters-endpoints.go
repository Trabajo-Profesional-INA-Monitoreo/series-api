package endpoints

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/controllers"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func setFiltersEndpoints(apiGroup *gin.RouterGroup, filterRepository repositories.FilterRepository) {
	controller := controllers.NewFilterController(services.NewFilterService(filterRepository))
	testApi := apiGroup.Group("/filtro")
	{
		testApi.GET("/procedimientos", controller.GetProcedures)
		testApi.GET("/estaciones", controller.GetStations)
		testApi.GET("/variables", controller.GetVariables)
		testApi.GET("/nodos", controller.GetNodes)
	}
	log.Infof("Configured filter endpoints")
}
