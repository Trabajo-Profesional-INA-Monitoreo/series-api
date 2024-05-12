package controllers

import (
	"errors"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	exceptions "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/errors"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SeriesController interface {
	GetStreamDataById(ctx *gin.Context)
	GetStreamCards(ctx *gin.Context)
	GetRedundancies(ctx *gin.Context)
}

type seriesController struct {
	seriesService services.StreamRetrievalService
}

func NewSeriesController(seriesService services.StreamService) SeriesController {
	return &seriesController{seriesService}
}

// GetStreamDataById godoc
//
//		@Summary		Endpoint para obtener los datos de una serie dado un id y su configuracion
//		@Tags           Series
//		@Produce		json
//		@Param          configuredStreamId      query     string  true  "Id de la configuracion de la serie"  Format(string)
//		@Param          timeStart    query     string  false  "Fecha de comienzo del periodo - valor por defecto: 7 dias atras"  Format(2006-01-02)
//		@Param          timeEnd      query     string  false  "Fecha del final del periodo - valor por defecto: mañana"  Format(2006-01-02)
//	    @Param          serie_id     path      int     true  "Id de la serie"
//		@Success		200	{object} dtos.StreamData
//		@Failure        400  {object}  dtos.ErrorResponse
//		@Failure        404  {object}  dtos.ErrorResponse
//		@Failure        500  {object}  dtos.ErrorResponse
//		@Router			/series/{serie_id} [get]
func (s seriesController) GetStreamDataById(ctx *gin.Context) {
	streamId, done := getUintPathParam(ctx, "serie_id")
	if done {
		return
	}
	configId, done := getUintQueryParam(ctx, "configuredStreamId")
	if done {
		return
	}
	timeStart, timeEnd, done := getDates(ctx)
	if done {
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
//	@Tags           Series
//	@Produce		json
//	@Param          timeStart    query     string  false  "Fecha de comienzo del periodo - valor por defecto: 7 dias atras"  Format(2006-01-02)
//	@Param          timeEnd      query     string  false  "Fecha del final del periodo - valor por defecto: mañana"  Format(2006-01-02)
//	@Param          streamId     query      int     false  "Filtro por ID de la serie"
//	@Param          stationId    query      int     false  "Filtro por ID de la estacion"
//	@Param          procId    	 query      int     false  "Filtro por ID de procedimiento"
//	@Param          varId    	 query      int     false  "Filtro por ID de variable"
//	@Param          nodeId    	 query      int     false  "Filtro por ID de nodo"
//	@Param          page    	 query      int     false  "Numero de pagina, por defecto 1"
//	@Param          pageSize     query      int     false  "Cantidad de series por pagina, por defecto 15"
//	@Param          configurationId     query      int     true  "ID de la configuracion"
//	@Success		200	{object} dtos.StreamCardsResponse
//	@Failure        400  {object}  dtos.ErrorResponse
//	@Failure        500  {object}  dtos.ErrorResponse
//	@Router			/series [get]
func (s seriesController) GetStreamCards(ctx *gin.Context) {
	timeStart, timeEnd, done := getDates(ctx)
	if done {
		return
	}
	configurationId, done := getUintQueryParam(ctx, "configurationId")
	if done {
		return
	}

	parameters := dtos.NewQueryParameters()
	parameters.AddParam("timeStart", timeStart)
	parameters.AddParam("timeEnd", timeEnd)
	parameters.AddParam("configurationId", configurationId)

	query, found := ctx.GetQuery("streamId")
	parameters.AddParamIfFound("streamId", query, found)

	query, found = ctx.GetQuery("nodeId")
	parameters.AddParamIfFound("nodeId", query, found)

	query, found = ctx.GetQuery("stationId")
	parameters.AddParamIfFound("stationId", query, found)

	query, found = ctx.GetQuery("stationId")
	parameters.AddParamIfFound("stationId", query, found)

	query, found = ctx.GetQuery("procId")
	parameters.AddParamIfFound("procId", query, found)

	query, found = ctx.GetQuery("varId")
	parameters.AddParamIfFound("varId", query, found)

	query = ctx.DefaultQuery("page", "1")
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

// GetRedundancies godoc
//
//	@Summary		Endpoint para obtener los ids de las redundancias de una serie configurada por id
//	@Tags           Series
//	@Produce		json
//	@Param          configured_stream_id     path      int     true  "Id de la serie configurada"
//	@Success		200	{object} dtos.Redundancies
//	@Failure        400  {object}  dtos.ErrorResponse
//	@Router			/series/redundancias/{configured_stream_id} [get]
func (s seriesController) GetRedundancies(ctx *gin.Context) {
	id, done := getUintPathParam(ctx, "configured_stream_id")
	if done {
		return
	}

	res := s.seriesService.GetRedundancies(id)

	ctx.JSON(http.StatusOK, res)
}
