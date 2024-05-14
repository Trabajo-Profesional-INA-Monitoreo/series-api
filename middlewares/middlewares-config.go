package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func SetUpMiddlewares(server *gin.Engine) {
	log.Infof("Setting up middlewares")
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	server.Use(gin.Recovery(), gin.Logger(), cors.New(config))
	setUpSwagger(server)

}
