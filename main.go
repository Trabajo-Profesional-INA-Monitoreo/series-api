package main

import (
	"fmt"
	"github.com/Trabajo-Profesional-INA-Monitoreo/inputs-api/config"
	"github.com/Trabajo-Profesional-INA-Monitoreo/inputs-api/endpoints"
	"github.com/Trabajo-Profesional-INA-Monitoreo/inputs-api/middlewares"
	"github.com/gin-gonic/gin"
	"log"
)

// @title			Inputs API
// @version		1.0
// @description	This API manages the inputs of the forecast model
func main() {
	apiConfig := config.GetConfig()
	server := gin.New()

	middlewares.SetUpMiddlewares(server)
	endpoints.SetUpEndpoints(server)

	err := server.Run(fmt.Sprintf(":%v", apiConfig.ServerPort))
	if err != nil {
		log.Fatalf("%v", err)
	}
}
