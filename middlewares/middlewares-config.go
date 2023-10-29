package middlewares

import (
	"github.com/gin-gonic/gin"
)

func SetUpMiddlewares(server *gin.Engine) {
	server.Use(gin.Recovery(), gin.Logger())
	setUpSwagger(server)
}
