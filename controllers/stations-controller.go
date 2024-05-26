package controllers

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services/stations-service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type StationsController interface {
	GetStations(ctx *gin.Context)
}

type stationsControllerImpl struct {
	stationsService stations_service.StationsService
}

func NewStationsController(stationsService stations_service.StationsService) StationsController {
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
//	@Param          page    	 query      int     false  "Numero de pagina, por defecto 1"
//	@Param          pageSize     query      int     false  "Cantidad de series por pagina, por defecto 15"
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
	parameters := dtos.NewQueryParameters()
	parameters.AddParam("timeStart", timeStart)
	parameters.AddParam("timeEnd", timeEnd)
	parameters.AddParam("configurationId", configId)
	query := ctx.DefaultQuery("page", "1")
	parameters.AddParam("page", query)

	query = ctx.DefaultQuery("pageSize", "15")
	parameters.AddParam("pageSize", query)
	res := s.stationsService.GetStations(parameters)
	ctx.JSON(http.StatusOK, res)
}
