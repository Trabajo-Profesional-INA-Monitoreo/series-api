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

// GetNetworks godoc
//
//	@Summary		Endpoint para obtener el resumen de las series agrupado por red
//	@Produce		json
//	@Success		200	{object} dtos.StreamsPerNetworkResponse
//	@Router			/series/redes [get]
func (s seriesController) GetNetworks(ctx *gin.Context) {
	res := s.seriesService.GetNetworks()
	ctx.JSON(http.StatusOK, res)
}

// GetStations godoc
//
//	@Summary		Endpoint para obtener el resumen de las series agrupado por estacion
//	@Produce		json
//	@Success		200	{object} dtos.StreamsPerStationResponse
//	@Router			/series/estaciones [get]
func (s seriesController) GetStations(ctx *gin.Context) {
	res := s.seriesService.GetStations()
	ctx.JSON(http.StatusOK, res)
}
