package endpoints

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func setUpHealthCheck(r *gin.RouterGroup) {
	r.GET("/healthcheck", healthcheck)
}

// healthcheck godoc
//
//	@Summary		Show the status of the server.
//	@Description	get the status of the server.
//	@Accept			*/*
//	@Produce		json
//	@Success		200	{string}	Server is up and running
//	@Router			/healthcheck [get]
func healthcheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Server is up and running")
}
