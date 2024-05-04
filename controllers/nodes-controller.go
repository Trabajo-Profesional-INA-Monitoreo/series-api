package controllers

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services/nodes-service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NodesController interface {
	GetNodes(ctx *gin.Context)
}

type nodesControllerImpl struct {
	nodesService nodes_service.NodesService
}

func NewNodesController(nodesService nodes_service.NodesService) NodesController {
	return &nodesControllerImpl{nodesService: nodesService}
}

// GetNodes godoc
//
//	@Summary		Endpoint para obtener el resumen de las series agrupado por nodo
//	@Tags           Series
//	@Produce		json
//	@Param          timeStart    query     string  false  "Fecha de comienzo del periodo - valor por defecto: 7 dias atras"  Format(2006-01-02)
//	@Param          timeEnd      query     string  false  "Fecha del final del periodo - valor por defecto: 5 dias despues"  Format(2006-01-02)
//	@Param          configurationId     query      int     true  "ID de la configuracion"
//	@Success		200	{object} dtos.StreamsPerNodeResponse
//	@Router			/series/nodos [get]
func (s nodesControllerImpl) GetNodes(ctx *gin.Context) {
	timeStart, timeEnd, done := getDates(ctx)
	if done {
		return
	}
	configId, done := getUintQueryParam(ctx, "configurationId")
	if done {
		return
	}
	res := s.nodesService.GetNodes(timeStart, timeEnd, configId)
	ctx.JSON(http.StatusOK, res)
}
