package endpoints

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/middlewares"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// @BasePath /api/v1

func SetUpEndpoints(server *gin.Engine, repositories *config.Repositories, apiConfig *config.ApiConfig) {

	log.Infof("Setting up endpoints")
	api := server.Group("/api/v1")
	if apiConfig.SecurityEnabled {
		api.Use(middlewares.IsAValidToken(apiConfig))
	}
	setSeriesEndpoints(api, repositories.StreamsRepository, clients.NewInaApiClientImpl(apiConfig))
	setInputsEndpoints(api, repositories.StreamsRepository)
	setConfigurationEndpoints(api, repositories.ConfigurationRepository)
	setErrorEndpoints(api, repositories.ErrorsRepository)
	setUpHealthCheck(api)
}
