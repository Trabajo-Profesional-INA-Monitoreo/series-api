package controllers

import (
	"fmt"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ConfigurationController interface {
	GetAllConfigurations(ctx *gin.Context)
	GetConfigurationById(ctx *gin.Context)
	CreateConfiguration(ctx *gin.Context)
	ModifyConfiguration(ctx *gin.Context)
	DeleteConfiguration(ctx *gin.Context)
}

type configurationController struct {
	configurationService services.ConfigurationService
}

// GetAllConfigurations godoc
//
//		@Summary		Endpoint para obtener las configuraciones
//		@Tags           Configuracion
//		@Produce		json
//		@Success		200	{array} dtos.AllConfigurations
//	    @Failure        400 {object} dtos.ErrorResponse
//		@Router			/configuracion [get]
func (c configurationController) GetAllConfigurations(ctx *gin.Context) {
	res := c.configurationService.GetAllConfigurations()
	ctx.JSON(http.StatusOK, res)
}

// GetConfigurationById godoc
//
//		@Summary		Endpoint para obtener una configuracion por id
//		@Tags           Configuracion
//		@Produce		json
//	    @Param          id     path      string     true  "Id de la configuracion"
//		@Success		200	{object} dtos.Configuration
//	    @Failure        400 {object} dtos.ErrorResponse
//		@Router			/configuracion/{id} [get]
func (c configurationController) GetConfigurationById(ctx *gin.Context) {
	id, done := getUintPathParam(ctx, "id")
	if done {
		return
	}
	res := c.configurationService.GetConfigurationById(id)

	if res == nil {
		ctx.JSON(http.StatusNotFound, dtos.NewErrorResponse(fmt.Errorf("Configuration not found")))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// CreateConfiguration godoc
//
//			@Summary		Endpoint para crear una configuracion
//			@Tags           Configuracion
//			@Produce		json
//	   		@Param          configuration  body  dtos.CreateConfiguration    true    "Add configuration"
//			@Success		201
//		    @Failure        400 {object} dtos.ErrorResponse
//		    @Failure        409 {object} dtos.ErrorResponse
//			@Router			/configuracion [post]
func (c configurationController) CreateConfiguration(ctx *gin.Context) {

	var configuration dtos.CreateConfiguration

	if err := ctx.ShouldBindJSON(&configuration); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(err))
		return
	}

	err := c.configurationService.CreateConfiguration(configuration)

	if err != nil {
		ctx.JSON(http.StatusConflict, dtos.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, "Configuration created")
}

// ModifyConfiguration godoc
//
//		@Summary		Endpoint para modificar una configuracion
//		@Tags           Configuracion
//		@Produce		json
//	   	@Param          configuration  body  dtos.Configuration    true    "Modify configuration"
//		@Success		200
//	    @Failure        400 {object} dtos.ErrorResponse
//	    @Failure        409 {object} dtos.ErrorResponse
//		@Router			/configuracion [put]
func (c configurationController) ModifyConfiguration(ctx *gin.Context) {
	var configuration dtos.Configuration

	if err := ctx.ShouldBindJSON(&configuration); err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(err))
		return
	}

	err := c.configurationService.ModifyConfiguration(configuration)

	if err != nil {
		ctx.JSON(http.StatusConflict, dtos.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Configuration modified")
}

// DeleteConfiguration godoc
//
//		@Summary		Endpoint para eliminar una configuracion por id
//		@Tags           Configuracion
//		@Produce		json
//	    @Param          id     path      string     true  "Id de la configuracion"
//		@Success		204
//	    @Failure        400 {object} dtos.ErrorResponse
//		@Router			/configuracion/{id} [delete]
func (c configurationController) DeleteConfiguration(ctx *gin.Context) {
	id, userSentId := ctx.Params.Get("id")
	if !userSentId {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("Id was not send")))
		return
	}
	c.configurationService.DeleteConfiguration(id)

	ctx.JSON(http.StatusOK, "Configuration deleted")
}

func NewConfigurationController(configurationService services.ConfigurationService) ConfigurationController {
	return &configurationController{configurationService}
}
