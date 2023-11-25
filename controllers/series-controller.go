package controllers

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SeriesController interface {
	GetNetworks(ctx *gin.Context)
	GetStations(ctx *gin.Context)
}

type seriesController struct {
	seriesService services.StreamService
}

func NewSeriesController(seriesService services.StreamService) SeriesController {
	return &seriesController{seriesService}
}

func (s seriesController) GetNetworks(ctx *gin.Context) {
	res := s.seriesService.GetNetworks()
	ctx.JSON(http.StatusOK, res)
}

func (s seriesController) GetStations(ctx *gin.Context) {

}
