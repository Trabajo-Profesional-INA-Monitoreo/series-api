package controllers

import (
	"fmt"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const DaysDefaultCured = 5
const DaysDefaultObservated = 1

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
//	@Param          timeStart    query     string  false  "Fecha de comienzo del periodo - valor por defecto: 7 dias atras"  Format(2006-01-02)
//	@Param          timeEnd      query     string  false  "Fecha del final del periodo - valor por defecto: 5 dias despues"  Format(2006-01-02)
//	@Success		200	{object} dtos.StreamsDataResponse
//	@Failure        400  {object}  dtos.ErrorResponse
//	@Router			/series/curadas/{serie_id} [get]
func (s seriesController) GetCuredSerieById(ctx *gin.Context) {
	id, userSentId := ctx.Params.Get("serie_id")
	if !userSentId {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("Id was not send")))
		return
	}

	timeStartQuery := ctx.DefaultQuery("timeStart", time.Now().Add(-DaysPerWeek*HoursPerDay*time.Hour).Format(time.DateOnly))
	timeEndQuery := ctx.DefaultQuery("timeEnd", time.Now().Add(DaysDefaultCured*HoursPerDay*time.Hour).Format(time.DateOnly))
	timeStart, err := time.Parse(time.DateOnly, timeStartQuery)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("error parsing time: %v", err)))
		return
	}
	timeEnd, err := time.Parse(time.DateOnly, timeEndQuery)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("error parsing time: %v", err)))
		return
	}

	res := s.seriesService.GetCuredSerieById(id, timeStart, timeEnd)

	ctx.JSON(http.StatusOK, res)
}

// GetObservatedSerieById godoc
//
//	@Summary		Endpoint para obtener los valores de una serie observada por id
//	@Produce		json
//	@Param          timeStart    query     string  false  "Fecha de comienzo del periodo - valor por defecto: 7 dias atras"  Format(2006-01-02)
//	@Param          timeEnd      query     string  false  "Fecha del final del periodo - valor por defecto: mañana"  Format(2006-01-02)
//	@Success		200	{object} dtos.StreamsDataResponse
//	@Failure        400  {object}  dtos.ErrorResponse
//	@Router			/series/observadas/{serie_id} [get]
func (s seriesController) GetObservatedSerieById(ctx *gin.Context) {
	id, userSentId := ctx.Params.Get("serie_id")
	if !userSentId {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("Id was not send")))
		return
	}
	timeStartQuery := ctx.DefaultQuery("timeStart", time.Now().Add(-DaysPerWeek*HoursPerDay*time.Hour).Format(time.DateOnly))
	timeEndQuery := ctx.DefaultQuery("timeEnd", time.Now().Add(DaysDefaultObservated*HoursPerDay*time.Hour).Format(time.DateOnly))
	timeStart, err := time.Parse(time.DateOnly, timeStartQuery)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("error parsing time: %v", err)))
		return
	}
	timeEnd, err := time.Parse(time.DateOnly, timeEndQuery)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("error parsing time: %v", err)))
		return
	}

	res := s.seriesService.GetObservatedSerieById(id, timeStart, timeEnd)

	ctx.JSON(http.StatusOK, res)
}

// GetPredictedSerieById godoc
//
//	@Summary		Endpoint para obtener los valores de una serie pronosticadas por id
//	@Produce		json
//	@Success		200	{object} dtos.CalibratedStreamsDataResponse
//	@Failure        400  {object}  dtos.ErrorResponse
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
