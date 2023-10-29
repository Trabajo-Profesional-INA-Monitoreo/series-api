package main

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/inputs-api/endpoints"
	"github.com/Trabajo-Profesional-INA-Monitoreo/inputs-api/middlewares"
	"github.com/gin-gonic/gin"
	"log"
)

//	@title			Inputs API
//	@version		1.0
//	@description	This API manages the inputs of the forecast model
func main() {
	server := gin.New()

	middlewares.SetUpMiddlewares(server)
	endpoints.SetUpEndpoints(server)

	err := server.Run("localhost:8080")
	if err != nil {
		log.Fatalf("%v", err)
	}
}
