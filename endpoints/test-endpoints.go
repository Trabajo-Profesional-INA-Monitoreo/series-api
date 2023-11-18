package endpoints

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/controllers"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
)

func setTestEndpoints(apiGroup *gin.RouterGroup) {
	controller := controllers.NewTestController(services.NewTestService())
	testApi := apiGroup.Group("/test")
	{
		testApi.GET("", controller.GetTest)
		testApi.POST("", controller.PostTest)
	}
}
