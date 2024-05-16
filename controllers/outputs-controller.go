package controllers

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services/outputs-service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OutputsController interface {
	GetOutputMetrics(ctx *gin.Context)
}

type outputsControllerImpl struct {
	outputsService outputs_service.OutputsService
}

func NewOutputsController(outputsService outputs_service.OutputsService) OutputsController {
	return &outputsControllerImpl{outputsService: outputsService}
}

// GetOutputMetrics godoc
//
//	@Summary		Endpoint para obtener las metricas de comportamiento
//	@Tags           Outputs
//	@Produce		json
//	@Param          timeStart    query     string  false  "Fecha de comienzo del periodo - valor por defecto: 7 dias atras"  Format(2006-01-02)
//	@Param          timeEnd      query     string  false  "Fecha del final del periodo - valor por defecto: mañana"  Format(2006-01-02)
//	@Param          configurationId     query      int     true  "ID de la configuracion"
//	@Success		200	{object} dtos.BehaviourStreamsResponse
//	@Failure        400  {object}  dtos.ErrorResponse
//	@Failure        500  {object}  dtos.ErrorResponse
//	@Router			/series/comportamiento [get]
func (s outputsControllerImpl) GetOutputMetrics(ctx *gin.Context) {
	timeStart, timeEnd, done := getDates(ctx)
	if done {
		return
	}
	configurationId, done := getUintQueryParam(ctx, "configurationId")
	if done {
		return
	}

	res, err := s.outputsService.GetOutputBehaviourMetrics(configurationId, timeStart, timeEnd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dtos.NewErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, res)
}
