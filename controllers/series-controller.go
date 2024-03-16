package controllers

import (
	"fmt"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SeriesController interface {
	GetNetworks(ctx *gin.Context)
	GetStations(ctx *gin.Context)
	GetCuredSerieById(ctx *gin.Context)
	GetObservatedSerieById(ctx *gin.Context)
	GetPredictedSerieById(ctx *gin.Context)
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

// GetCuredSerieById godoc
//
//	@Summary		Endpoint para obtener los valores de una serie curada por id
//	@Produce		json
//	@Success		200	{object} dtos.StreamsDataResponse
//	@Router			/series/curadas/{serie_id} [get]
func (s seriesController) GetCuredSerieById(ctx *gin.Context) {
	id, userSentId := ctx.Params.Get("serie_id")
	if !userSentId {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("Id was not send")))
		return
	}
	res := s.seriesService.GetCuredSerieById(id)

	ctx.JSON(http.StatusOK, res)
}

// GetObservatedSerieById godoc
//
//	@Summary		Endpoint para obtener los valores de una serie observada por id
//	@Produce		json
//	@Success		200	{object} dtos.StreamsDataResponse
//	@Router			/series/observadas/{serie_id} [get]
func (s seriesController) GetObservatedSerieById(ctx *gin.Context) {
	id, userSentId := ctx.Params.Get("serie_id")
	if !userSentId {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("Id was not send")))
		return
	}
	res := s.seriesService.GetObservatedSerieById(id)

	ctx.JSON(http.StatusOK, res)
}

// GetPredictedSerieById godoc
//
//	@Summary		Endpoint para obtener los valores de una serie pronosticadas por id
//	@Produce		json
//	@Success		200	{object} dtos.CalibratedStreamsDataResponse
//	@Router			/series/pronosticadas/{calibrado_id} [get]
func (s seriesController) GetPredictedSerieById(ctx *gin.Context) {
	id, userSentId := ctx.Params.Get("calibrado_id")
	if !userSentId {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("Id was not send")))
		return
	}
	res := s.seriesService.GetPredictedSerieById(id)

	ctx.JSON(http.StatusOK, res)
}
