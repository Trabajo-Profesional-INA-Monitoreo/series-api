package controllers

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InputsController interface {
	GetGeneralMetrics(ctx *gin.Context)
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
//	@Produce		json
//	@Success		200	{object} dtos.InputsGeneralMetrics
//	@Router			/inputs/metricas-generales [get]
func (i inputsController) GetGeneralMetrics(ctx *gin.Context) {
	res := i.inputsService.GetGeneralMetrics()
	ctx.JSON(http.StatusOK, res)
}
