package endpoints

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/middlewares"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// @BasePath /api/v1

func SetUpEndpoints(server *gin.Engine, repositories *config.Repositories, apiConfig *config.ServiceConfigurationData) {

	log.Infof("Setting up endpoints")
	api := server.Group("/api/v1")
	if apiConfig.SecurityEnabled {
		api.Use(middlewares.IsAValidToken(apiConfig))
	}
	inaClient := clients.NewInaApiClientImpl(apiConfig)
	setSeriesEndpoints(api, repositories, inaClient)
	setInputsEndpoints(api, repositories.InputsRepository)
	setConfigurationEndpoints(api, repositories, inaClient, apiConfig)
	setErrorEndpoints(api, repositories.ErrorsRepository)
	setFiltersEndpoints(api, repositories.FilterRepository)
	SetNotificationEndpoints(api, apiConfig)
	setUpHealthCheck(api)
}
