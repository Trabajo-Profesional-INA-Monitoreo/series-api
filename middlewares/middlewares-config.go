package middlewares

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func SetUpMiddlewares(server *gin.Engine) {
	log.Infof("Setting up middlewares")
	server.Use(gin.Recovery(), gin.Logger())
	setUpSwagger(server)
}
