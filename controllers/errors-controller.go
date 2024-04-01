package controllers

import (
	"fmt"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
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
//	@Success		200	 {array}   dtos.ErrorsCountPerDayAndType
//	@Failure        400  {object}  dtos.ErrorResponse
//	@Router			/errores/por-dia [get]
func (e errorsController) GetErrorsPerDay(ctx *gin.Context) {
	timeStart, timeEnd, done := getDates(ctx)
	if done {
		return
	}
	response := e.errorsService.GetErrorsPerDay(timeStart, timeEnd)
	ctx.JSON(http.StatusOK, response)
}

func getDates(ctx *gin.Context) (time.Time, time.Time, bool) {
	timeStartQuery := ctx.DefaultQuery("timeStart", time.Now().Add(-DaysPerWeek*HoursPerDay*time.Hour).Format(time.DateOnly))
	timeEndQuery := ctx.DefaultQuery("timeEnd", time.Now().Add(HoursPerDay*time.Hour).Format(time.DateOnly))
	timeStart, err := time.Parse(time.DateOnly, timeStartQuery)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("error parsing time: %v", err)))
		return time.Time{}, time.Time{}, true
	}
	timeEnd, err := time.Parse(time.DateOnly, timeEndQuery)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("error parsing time: %v", err)))
		return time.Time{}, time.Time{}, true
	}
	return timeStart, timeEnd, false
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
	configurationIdQuery, sent := ctx.GetQuery("configurationId")
	if !sent {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("configurationId missing")))
		return
	}
	configurationId, err := strconv.ParseUint(configurationIdQuery, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("wrong format for configurationId, should be uint")))
		return
	}
	result := e.errorsService.GetErrorIndicators(timeStart, timeEnd, configurationId)
	ctx.JSON(http.StatusOK, result)
}

func NewErrorsController(errorsService services.ErrorsService) ErrorsController {
	return &errorsController{errorsService: errorsService}
}
