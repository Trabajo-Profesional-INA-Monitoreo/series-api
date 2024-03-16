package controllers

import (
	"fmt"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const HoursPerDay = 24
const DaysPerWeek = 7

type ErrorsController interface {
	GetErrorsPerDay(ctx *gin.Context)
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
	timeStartQuery := ctx.DefaultQuery("timeStart", time.Now().Add(-DaysPerWeek*HoursPerDay*time.Hour).Format(time.DateOnly))
	timeEndQuery := ctx.DefaultQuery("timeEnd", time.Now().Add(HoursPerDay*time.Hour).Format(time.DateOnly))
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
	response := e.errorsService.GetErrorsPerDay(timeStart, timeEnd)
	ctx.JSON(http.StatusOK, response)
}

func NewErrorsController(errorsService services.ErrorsService) ErrorsController {
	return &errorsController{errorsService: errorsService}
}
