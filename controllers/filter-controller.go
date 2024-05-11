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
	GetNodes(ctx *gin.Context)
}

type filterController struct {
	filterService services.FilterService
}

// GetProcedures godoc
//
//		@Summary		Endpoint para obtener los procedimientos
//		@Tags           Filtros
//		@Produce		json
//	@Param          configurationId     query      int     true  "ID de la configuracion"
//		@Success		200	{array} dtos.FilterValue
//	    @Failure        400 {object} dtos.ErrorResponse
//		@Router			/filtro/procedimientos [get]
func (f filterController) GetProcedures(ctx *gin.Context) {
	configurationId, done := getUintQueryParam(ctx, "configurationId")
	if done {
		return
	}
	res := f.filterService.GetProcedures(configurationId)
	ctx.JSON(http.StatusOK, res)
}

// GetStations godoc
//
//		@Summary		Endpoint para obtener las estaciones
//		@Tags           Filtros
//		@Produce		json
//	@Param          configurationId     query      int     true  "ID de la configuracion"
//		@Success		200	{array} dtos.FilterValue
//	    @Failure        400 {object} dtos.ErrorResponse
//		@Router			/filtro/estaciones [get]
func (f filterController) GetStations(ctx *gin.Context) {
	configurationId, done := getUintQueryParam(ctx, "configurationId")
	if done {
		return
	}
	res := f.filterService.GetStations(configurationId)
	ctx.JSON(http.StatusOK, res)
}

// GetVariables godoc
//
//		@Summary		Endpoint para obtener las variables
//		@Tags           Filtros
//		@Produce		json
//	@Param          configurationId     query      int     true  "ID de la configuracion"
//		@Success		200	{array} dtos.FilterValue
//	    @Failure        400 {object} dtos.ErrorResponse
//		@Router			/filtro/variables [get]
func (f filterController) GetVariables(ctx *gin.Context) {
	configurationId, done := getUintQueryParam(ctx, "configurationId")
	if done {
		return
	}
	res := f.filterService.GetVariables(configurationId)
	ctx.JSON(http.StatusOK, res)
}

// GetNodes godoc
//
//		@Summary		Endpoint para obtener los nodos
//		@Tags           Filtros
//		@Produce		json
//	@Param          configurationId     query      int     true  "ID de la configuracion"
//		@Success		200	{array} dtos.FilterValue
//	    @Failure        400 {object} dtos.ErrorResponse
//		@Router			/filtro/nodos [get]
func (f filterController) GetNodes(ctx *gin.Context) {
	configurationId, done := getUintQueryParam(ctx, "configurationId")
	if done {
		return
	}
	res := f.filterService.GetNodes(configurationId)
	ctx.JSON(http.StatusOK, res)
}

func NewFilterController(service services.FilterService) FilterController {
	return &filterController{filterService: service}
}
