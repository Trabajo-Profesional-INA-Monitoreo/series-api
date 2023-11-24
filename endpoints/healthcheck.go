package endpoints

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func setUpHealthCheck(r *gin.RouterGroup) {
	r.GET("/healthcheck", healthcheck)
	log.Infof("Configured healthcheck endpoint")
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
