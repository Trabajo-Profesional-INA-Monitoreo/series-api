package controllers

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type StationsController interface {
	GetStations(ctx *gin.Context)
}

type stationsControllerImpl struct {
	stationsService services.StationsService
}

func NewStationsController(stationsService services.StationsService) StationsController {
	return &stationsControllerImpl{stationsService: stationsService}
}

// GetStations godoc
//
//	@Summary		Endpoint para obtener el resumen de las series agrupado por estacion
//	@Tags           Series
//	@Produce		json
//	@Param          timeStart    query     string  false  "Fecha de comienzo del periodo - valor por defecto: 7 dias atras"  Format(2006-01-02)
//	@Param          timeEnd      query     string  false  "Fecha del final del periodo - valor por defecto: 5 dias despues"  Format(2006-01-02)
//	@Param          configurationId     query      int     true  "ID de la configuracion"
//	@Success		200	{object} dtos.StreamsPerStationResponse
//	@Router			/series/estaciones [get]
func (s stationsControllerImpl) GetStations(ctx *gin.Context) {
	timeStart, timeEnd, done := getDates(ctx)
	if done {
		return
	}
	configId, done := getUintQueryParam(ctx, "configurationId")
	if done {
		return
	}
	res := s.stationsService.GetStations(timeStart, timeEnd, configId)
	ctx.JSON(http.StatusOK, res)
}
