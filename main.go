package main

import (
	"fmt"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/config"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/endpoints"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/jobs"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/middlewares"
	"github.com/gin-gonic/gin"
	"log"
)

// @title			Inputs API
// @version		1.0
// @description	This API manages the inputs of the forecast model
func main() {
	apiConfig := config.GetConfig()
	config.InitLogger(apiConfig.LogLevel)
	server := gin.New()

	repositories := config.CreateRepositories(apiConfig.DbUrl, apiConfig.SqlLogLevel)
	middlewares.SetUpMiddlewares(server)
	endpoints.SetUpEndpoints(server, repositories, apiConfig)
	jobs.SetUpJobs(repositories, apiConfig)

	err := server.Run(fmt.Sprintf(":%v", apiConfig.ServerPort))
	if err != nil {
		log.Fatalf("%v", err)
	}
}
