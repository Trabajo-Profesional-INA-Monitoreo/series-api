package endpoints

import "github.com/gin-gonic/gin"

// @BasePath /api/v1

func SetUpEndpoints(server *gin.Engine) {
	api := server.Group("/api/v1")
	setTestEndpoints(api)
	setUpHealthCheck(api)
}
