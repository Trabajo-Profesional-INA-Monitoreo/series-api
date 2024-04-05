package controllers

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

const HoursPerDay = 24
const DaysPerWeek = 7

type ErrorsController interface {
	GetErrorsPerDay(ctx *gin.Context)
	GetErrorIndicators(context *gin.Context)
}

type errorsController struct {
	errorsService services.ErrorsService
}

// GetErrorsPerDay godoc
//
//	@Summary		Endpoint para obtener las errores detectados por dia
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
	configurationId, done := getConfigurationId(ctx)
	if done {
		return
	}
	response := e.errorsService.GetErrorsPerDay(timeStart, timeEnd, configurationId)
	ctx.JSON(http.StatusOK, response)
}

// GetErrorIndicators godoc
//
//	@Summary		Endpoint para obtener las indicadores de errores
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
	configurationId, done := getConfigurationId(ctx)
	if done {
		return
	}
	result := e.errorsService.GetErrorIndicators(timeStart, timeEnd, configurationId)
	ctx.JSON(http.StatusOK, result)
}

func NewErrorsController(errorsService services.ErrorsService) ErrorsController {
	return &errorsController{errorsService: errorsService}
}
