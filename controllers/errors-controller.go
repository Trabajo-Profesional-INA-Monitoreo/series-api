package controllers

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

const HoursPerDay = 24
const DaysPerWeek = 7

type ErrorsController interface {
	GetErrorsPerDay(ctx *gin.Context)
	GetErrorIndicators(context *gin.Context)
	GetStreamsWithRelatedError(ctx *gin.Context)
	GetErrorsOfConfiguredStream(ctx *gin.Context)
}

type errorsController struct {
	errorsService services.ErrorsService
}

// GetErrorsPerDay godoc
//
//	@Summary		Endpoint para obtener las errores detectados por dia
//	@Tags           Errores
//	@Produce		json
//	@Param          timeStart    query     string  false  "Fecha de comienzo del periodo - valor por defecto: 7 dias atras"  Format(2006-01-02)
//	@Param          timeEnd      query     string  false  "Fecha del final del periodo - valor por defecto: hoy"  Format(2006-01-02)
//	@Param          configurationId     query      int     true  "ID de la configuracion"
//	@Success		200	 {array}   dtos.ErrorsCountPerDayAndType
//	@Failure        400  {object}  dtos.ErrorResponse
//	@Router			/errores/por-dia [get]
func (e errorsController) GetErrorsPerDay(ctx *gin.Context) {
	timeStart, timeEnd, done := getDates(ctx)
	if done {
		return
	}
	configurationId, done := getUintQueryParam(ctx, "configurationId")
	if done {
		return
	}
	response := e.errorsService.GetErrorsPerDay(timeStart, timeEnd, configurationId)
	ctx.JSON(http.StatusOK, response)
}

// GetErrorIndicators godoc
//
//	@Summary		Endpoint para obtener las indicadores de errores
//	@Tags           Errores
//	@Produce		json
//	@Param          timeStart    query     string  false  "Fecha de comienzo del periodo - valor por defecto: 7 dias atras"  Format(2006-01-02)
//	@Param          timeEnd      query     string  false  "Fecha del final del periodo - valor por defecto: hoy"  Format(2006-01-02)
//	@Param          configurationId      query     string  true  "Id de la configuracion"  Format(uint)
//	@Success		200	 {array}   dtos.ErrorIndicator
//	@Failure        400  {object}  dtos.ErrorResponse
//	@Router			/errores/indicadores [get]
func (e errorsController) GetErrorIndicators(ctx *gin.Context) {
	timeStart, timeEnd, done := getDates(ctx)
	if done {
		return
	}
	configurationId, done := getUintQueryParam(ctx, "configurationId")
	if done {
		return
	}
	result := e.errorsService.GetErrorIndicators(timeStart, timeEnd, configurationId)
	ctx.JSON(http.StatusOK, result)
}

// GetStreamsWithRelatedError godoc
//
//	@Summary		Endpoint para obtener las indicadores de errores
//	@Tags           Errores
//	@Produce		json
//	@Param          timeStart    query     string  false  "Fecha de comienzo del periodo - valor por defecto: 7 dias atras"  Format(2006-01-02)
//	@Param          timeEnd      query     string  false  "Fecha del final del periodo - valor por defecto: hoy"  Format(2006-01-02)
//	@Param          configurationId      query     string  true  "Id de la configuracion"  Format(uint)
//	@Param          errorType      query     string  true  "Id del tipo de error"  Format(uint)
//	@Success		200	 {array}   dtos.ErrorIndicator
//	@Failure        400  {object}  dtos.ErrorResponse
//	@Router			/errores/series-implicadas [get]
func (e errorsController) GetStreamsWithRelatedError(ctx *gin.Context) {
	timeStart, timeEnd, done := getDates(ctx)
	if done {
		return
	}
	configurationId, done := getUintQueryParam(ctx, "configurationId")
	if done {
		return
	}
	errorId, done := getUintQueryParam(ctx, "errorType")
	if done {
		return
	}
	parameters := dtos.NewQueryParameters()
	parameters.AddParam("configurationId", configurationId)
	parameters.AddParam("timeStart", timeStart)
	parameters.AddParam("timeEnd", timeEnd)
	parameters.AddParam("errorType", errorId)
	result, err := e.errorsService.GetRelatedStreams(parameters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.NewErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// GetErrorsOfConfiguredStream godoc
//
//	@Summary		Endpoint para obtener los errores de una serie dado un id
//	@Tags           Errores
//	@Produce		json
//	@Param          configuredStreamId      path     int  true  "Id de la configuracion de la serie"  Format(int)
//	@Param          timeStart    query     string  false  "Fecha de comienzo del periodo - valor por defecto: 7 dias atras"  Format(2006-01-02)
//	@Param          timeEnd      query     string  false  "Fecha del final del periodo - valor por defecto: mañana"  Format(2006-01-02)
//	@Param          page    	 query      int     false  "Numero de pagina, por defecto 1"
//	@Param          pageSize     query      int     false  "Cantidad de series por pagina, por defecto 15"
//	@Success		200	{object} 	dtos.DetectedErrorsOfStream
//	@Failure        400  {object}  dtos.ErrorResponse
//	@Failure        500  {object}  dtos.ErrorResponse
//	@Router			/errores/{configuredStreamId} [get]
func (e errorsController) GetErrorsOfConfiguredStream(ctx *gin.Context) {
	configStreamId, done := getUintPathParam(ctx, "configuredStreamId")
	if done {
		return
	}
	timeStart, timeEnd, done := getDates(ctx)
	if done {
		return
	}
	parameters := dtos.NewQueryParameters()
	query := ctx.DefaultQuery("page", "1")
	parameters.AddParam("page", query)

	query = ctx.DefaultQuery("pageSize", "15")
	parameters.AddParam("pageSize", query)
	parameters.AddParam("timeStart", timeStart)
	parameters.AddParam("timeEnd", timeEnd)
	parameters.AddParam("configStreamId", configStreamId)
	res, err := e.errorsService.GetErrorsOfConfiguredStream(parameters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.NewErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func NewErrorsController(errorsService services.ErrorsService) ErrorsController {
	return &errorsController{errorsService: errorsService}
}
