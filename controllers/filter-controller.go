package controllers

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FilterController interface {
	GetProcedures(ctx *gin.Context)
	GetStations(ctx *gin.Context)
	GetVariables(ctx *gin.Context)
}

type filterController struct {
	filterService services.FilterService
}

// GetProcedures godoc
//
//		@Summary		Endpoint para obtener los procedimientos
//		@Tags           Filtros
//		@Produce		json
//		@Success		200	{array} dtos.ProcedureFilter
//	    @Failure        400 {object} dtos.ErrorResponse
//		@Router			/filtro/procedimientos [get]
func (f filterController) GetProcedures(ctx *gin.Context) {
	res := f.filterService.GetProcedures()
	ctx.JSON(http.StatusOK, res)
}

// GetStations godoc
//
//		@Summary		Endpoint para obtener las estaciones
//		@Tags           Filtros
//		@Produce		json
//		@Success		200	{array} dtos.StationFilter
//	    @Failure        400 {object} dtos.ErrorResponse
//		@Router			/filtro/estaciones [get]
func (f filterController) GetStations(ctx *gin.Context) {
	res := f.filterService.GetStations()
	ctx.JSON(http.StatusOK, res)
}

// GetVariables godoc
//
//		@Summary		Endpoint para obtener las variables
//		@Tags           Filtros
//		@Produce		json
//		@Success		200	{array} dtos.VariableFilter
//	    @Failure        400 {object} dtos.ErrorResponse
//		@Router			/filtro/variables [get]
func (f filterController) GetVariables(ctx *gin.Context) {
	res := f.filterService.GetVariables()
	ctx.JSON(http.StatusOK, res)
}

func NewFilterController(service services.FilterService) FilterController {
	return &filterController{filterService: service}
}
