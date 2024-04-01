package controllers

import (
	"errors"
	"fmt"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	exceptions "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/errors"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	GetStreamDataById(ctx *gin.Context)
	GetStreamCards(ctx *gin.Context)
	GetOutputMetrics(ctx *gin.Context)
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
//	@Param          serie_id     path      int     true  "ID de la serie"
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
//	@Param          serie_id     path      int     true  "ID de la serie"
//	@Success		200	{object} dtos.StreamsDataResponse
//	@Failure        400  {object}  dtos.ErrorResponse
//	@Router			/series/observadas/{serie_id} [get]
func (s seriesController) GetObservatedSerieById(ctx *gin.Context) {
	id, userSentId := ctx.Params.Get("serie_id")
	if !userSentId {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("serie_id was not send")))
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
//	@Param          calibrado_id     path      int     true  "ID del calibrado"
//	@Success		200	{object} dtos.CalibratedStreamsDataResponse
//	@Failure        400  {object}  dtos.ErrorResponse
//	@Router			/series/pronosticadas/{calibrado_id} [get]
func (s seriesController) GetPredictedSerieById(ctx *gin.Context) {
	id, userSentId := ctx.Params.Get("calibrado_id")
	if !userSentId {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("calibrado_id was not sent")))
		return
	}
	res := s.seriesService.GetPredictedSerieById(id)

	ctx.JSON(http.StatusOK, res)
}

// GetStreamDataById godoc
//
//		@Summary		Endpoint para obtener los datos de una serie dado un id y su configuracion
//		@Produce		json
//		@Param          configuredStreamId      query     string  true  "Id de la configuracion de la serie"  Format(string)
//		@Param          timeStart    query     string  false  "Fecha de comienzo del periodo - valor por defecto: 7 dias atras"  Format(2006-01-02)
//		@Param          timeEnd      query     string  false  "Fecha del final del periodo - valor por defecto: mañana"  Format(2006-01-02)
//	    @Param          serie_id     path      int     true  "ID de la serie"
//		@Success		200	{object} dtos.StreamData
//		@Failure        400  {object}  dtos.ErrorResponse
//		@Failure        404  {object}  dtos.ErrorResponse
//		@Failure        500  {object}  dtos.ErrorResponse
//		@Router			/series/{serie_id} [get]
func (s seriesController) GetStreamDataById(ctx *gin.Context) {
	streamIdParam, userSentId := ctx.Params.Get("serie_id")
	if !userSentId {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("serie_id was not sent")))
		return
	}
	configIdParam, userSentConfigId := ctx.GetQuery("configuredStreamId")
	if !userSentConfigId {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("configuredStreamId was not sent")))
		return
	}
	streamId, err := strconv.ParseUint(streamIdParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("serie_id should be a number")))
		return
	}
	configId, err := strconv.ParseUint(configIdParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("configuredStreamId should be a number")))
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
	streamData, err := s.seriesService.GetStreamData(streamId, configId, timeStart, timeEnd)
	if errors.Is(err, &exceptions.NotFound{}) {
		ctx.JSON(http.StatusNotFound, dtos.NewErrorResponse(err))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.NewErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, streamData)
}

// GetStreamCards godoc
//
//	@Summary		Endpoint para obtener las series configuradas de una configuracion
//	@Produce		json
//	@Param          timeStart    query     string  false  "Fecha de comienzo del periodo - valor por defecto: 7 dias atras"  Format(2006-01-02)
//	@Param          timeEnd      query     string  false  "Fecha del final del periodo - valor por defecto: mañana"  Format(2006-01-02)
//	@Param          streamId     query      int     false  "Filtro por ID de la serie"
//	@Param          stationId    query      int     false  "Filtro por ID de la estacion"
//	@Param          procId    	 query      int     false  "Filtro por ID de procedimiento"
//	@Param          varId    	 query      int     false  "Filtro por ID de variable"
//	@Param          page    	 query      int     false  "Numero de pagina, por defecto 0"
//	@Param          pageSize     query      int     false  "Cantidad de series por pagina, por defecto 15"
//	@Param          configurationId     query      int     true  "ID de la configuracion"
//	@Success		200	{object} dtos.StreamCardsResponse
//	@Failure        400  {object}  dtos.ErrorResponse
//	@Failure        500  {object}  dtos.ErrorResponse
//	@Router			/series [get]
func (s seriesController) GetStreamCards(ctx *gin.Context) {
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
	configId, found := ctx.GetQuery("configurationId")
	if !found {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("configurationId missing")))
		return
	}

	parameters := dtos.NewStreamCardsParameters()
	parameters.AddParam("timeStart", timeStart)
	parameters.AddParam("timeEnd", timeEnd)
	parameters.AddParam("configurationId", configId)

	query, found := ctx.GetQuery("streamId")
	parameters.AddParamIfFound("streamId", query, found)

	query, found = ctx.GetQuery("stationId")
	parameters.AddParamIfFound("stationId", query, found)

	query, found = ctx.GetQuery("stationId")
	parameters.AddParamIfFound("stationId", query, found)

	query, found = ctx.GetQuery("procId")
	parameters.AddParamIfFound("procId", query, found)

	query, found = ctx.GetQuery("varId")
	parameters.AddParamIfFound("varId", query, found)

	query = ctx.DefaultQuery("page", "0")
	parameters.AddParam("page", query)

	query = ctx.DefaultQuery("pageSize", "15")
	parameters.AddParam("pageSize", query)

	res, err := s.seriesService.GetStreamCards(parameters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// GetOutputMetrics godoc
//
//	@Summary		Endpoint para obtener las metricas de comportamiento
//	@Produce		json
//	@Param          timeStart    query     string  false  "Fecha de comienzo del periodo - valor por defecto: 7 dias atras"  Format(2006-01-02)
//	@Param          timeEnd      query     string  false  "Fecha del final del periodo - valor por defecto: mañana"  Format(2006-01-02)
//	@Param          configurationId     query      int     true  "ID de la configuracion"
//	@Success		200	{object} dtos.BehaviourStreamsResponse
//	@Failure        400  {object}  dtos.ErrorResponse
//	@Failure        500  {object}  dtos.ErrorResponse
//	@Router			/series/comportamiento [get]
func (s seriesController) GetOutputMetrics(ctx *gin.Context) {
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
	configIdQuery, found := ctx.GetQuery("configurationId")
	if !found {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("configurationId missing")))
		return
	}
	configId, err := strconv.ParseUint(configIdQuery, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("configurationId is not an int")))
		return
	}

	res, err := s.seriesService.GetOutputBehaviourMetrics(configId, timeStart, timeEnd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.NewErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, res)
}
