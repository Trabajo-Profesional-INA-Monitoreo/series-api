package controllers

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InputsController interface {
	GetGeneralMetrics(ctx *gin.Context)
	GetTotalStreamsWithNullValues(ctx *gin.Context)
	GetTotalStreamsWithObservedOutlier(context *gin.Context)
}

type inputsController struct {
	inputsService services.InputsService
}

func NewInputsController(inputsService services.InputsService) InputsController {
	return &inputsController{inputsService}
}

// GetGeneralMetrics godoc
//
//	@Summary		Endpoint para obtener las metricas generales de inputs
//	@Tags           Inputs
//	@Produce		json
//	@Param          configurationId     query      int     true  "ID de la configuracion"
//	@Success		200	{object} dtos.InputsGeneralMetrics
//	@Router			/inputs/metricas-generales [get]
func (i inputsController) GetGeneralMetrics(ctx *gin.Context) {
	configurationId, done := getUintQueryParam(ctx, "configurationId")
	if done {
		return
	}
	res := i.inputsService.GetGeneralMetrics(configurationId)
	ctx.JSON(http.StatusOK, res)
}

// GetTotalStreamsWithNullValues godoc
//
//	@Summary		Endpoint para obtener la cantidad de series con valores nulos
//	@Tags           Inputs
//	@Produce		json
//	@Param          timeStart    query     string  false  "Fecha de comienzo del periodo - valor por defecto: 7 dias atras"  Format(2006-01-02)
//	@Param          timeEnd      query     string  false  "Fecha del final del periodo - valor por defecto: mañana"  Format(2006-01-02)
//	@Param          configurationId     query      int     true  "ID de la configuracion"
//	@Success		200	{object} dtos.TotalStreamsWithNullValues
//	@Router			/inputs/series-con-nulos [get]
func (i inputsController) GetTotalStreamsWithNullValues(ctx *gin.Context) {
	configurationId, done := getUintQueryParam(ctx, "configurationId")
	if done {
		return
	}
	timeStart, timeEnd, done := getDates(ctx)
	if done {
		return
	}
	res := i.inputsService.GetTotalStreamsWithNullValues(configurationId, timeStart, timeEnd)
	ctx.JSON(http.StatusOK, res)
}

// GetTotalStreamsWithObservedOutlier godoc
//
//	@Summary		Endpoint para obtener la cantidad de series con valores fuera de los umbrales
//	@Tags           Inputs
//	@Produce		json
//	@Param          timeStart    query     string  false  "Fecha de comienzo del periodo - valor por defecto: 7 dias atras"  Format(2006-01-02)
//	@Param          timeEnd      query     string  false  "Fecha del final del periodo - valor por defecto: mañana"  Format(2006-01-02)
//	@Param          configurationId     query      int     true  "ID de la configuracion"
//	@Success		200	{object} dtos.TotalStreamsWithObservedOutlier
//	@Router			/inputs/series-fuera-umbral [get]
func (i inputsController) GetTotalStreamsWithObservedOutlier(ctx *gin.Context) {
	configurationId, done := getUintQueryParam(ctx, "configurationId")
	if done {
		return
	}
	timeStart, timeEnd, done := getDates(ctx)
	if done {
		return
	}
	res := i.inputsService.GetTotalStreamsWithObservedOutlier(configurationId, timeStart, timeEnd)
	ctx.JSON(http.StatusOK, res)
}
