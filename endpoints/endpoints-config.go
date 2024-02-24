package endpoints

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// @BasePath /api/v1

func SetUpEndpoints(server *gin.Engine, configArguments *config.ApiConfig) {
	repositories := config.CreateRepositories(configArguments.DbUrl)
	log.Infof("Setting up endpoints")
	api := server.Group("/api/v1")
	setSeriesEndpoints(api, repositories.StreamsRepository)
	setInputsEndpoints(api, repositories.StreamsRepository)
	setConfigurationEndpoints(api, repositories.ConfigurationRepository)
	setUpHealthCheck(api)
}
