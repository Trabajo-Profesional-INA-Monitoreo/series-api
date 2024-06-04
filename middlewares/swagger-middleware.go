package middlewares

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setUpSwagger(server *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/api/v1"
	server.GET("/api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
