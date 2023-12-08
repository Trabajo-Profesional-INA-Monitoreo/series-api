package endpoints

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// @BasePath /api/v1

func SetUpEndpoints(server *gin.Engine, config *config.ApiConfig) {
	log.Infof("Setting up endpoints")
	api := server.Group("/api/v1")
	setSeriesEndpoints(api, config)
	setInputsEndpoints(api, config)
	setUpHealthCheck(api)
}
