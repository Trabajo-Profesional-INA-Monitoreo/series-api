package controllers

import (
	"errors"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	exceptions "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/errors"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InaController interface {
	GetCuredSerieById(ctx *gin.Context)
	GetObservatedSerieById(ctx *gin.Context)
	GetPredictedSerieById(ctx *gin.Context)
}

type inaControllerImpl struct {
	inaService services.InaServiceApi
}

func NewInaController(inaService services.InaServiceApi) InaController {
	return &inaControllerImpl{inaService: inaService}
}

// GetCuredSerieById godoc
//
//	@Summary		Endpoint para obtener los valores de una serie curada por id
//	@Tags           Interfaz INA
//	@Produce		json
//	@Param          timeStart    query     string  false  "Fecha de comienzo del periodo - valor por defecto: 7 dias atras"  Format(2006-01-02)
//	@Param          timeEnd      query     string  false  "Fecha del final del periodo - valor por defecto: 5 dias despues"  Format(2006-01-02)
//	@Param          serie_id     path      int     true  "Id de la serie"
//	@Success		200	{object} dtos.StreamsDataResponse
//	@Failure        400  {object}  dtos.ErrorResponse
//	@Failure        404  {object}  dtos.ErrorResponse
//	@Router			/series/curadas/{serie_id} [get]
func (s inaControllerImpl) GetCuredSerieById(ctx *gin.Context) {
	id, done := getUintPathParam(ctx, "serie_id")
	if done {
		return
	}
	timeStart, timeEnd, done := getDates(ctx)
	if done {
		return
	}

	res, err := s.inaService.GetCuredSerieById(id, timeStart, timeEnd)
	if errors.Is(err, &exceptions.NotFound{}) {
		ctx.JSON(http.StatusNotFound, dtos.NewErrorResponse(err))
		return
	}
	if errors.Is(err, &exceptions.BadRequest{}) {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// GetObservatedSerieById godoc
//
//	@Summary		Endpoint para obtener los valores de una serie observada por id
//	@Tags           Interfaz INA
//	@Produce		json
//	@Param          timeStart    query     string  false  "Fecha de comienzo del periodo - valor por defecto: 7 dias atras"  Format(2006-01-02)
//	@Param          timeEnd      query     string  false  "Fecha del final del periodo - valor por defecto: mañana"  Format(2006-01-02)
//	@Param          serie_id     path      int     true  "Id de la serie"
//	@Success		200	{object} dtos.StreamsDataResponse
//	@Failure        400  {object}  dtos.ErrorResponse
//	@Failure        404  {object}  dtos.ErrorResponse
//	@Router			/series/observadas/{serie_id} [get]
func (s inaControllerImpl) GetObservatedSerieById(ctx *gin.Context) {
	id, done := getUintPathParam(ctx, "serie_id")
	if done {
		return
	}
	timeStart, timeEnd, done := getDates(ctx)
	if done {
		return
	}

	res, err := s.inaService.GetObservatedSerieById(id, timeStart, timeEnd)
	if errors.Is(err, &exceptions.NotFound{}) {
		ctx.JSON(http.StatusNotFound, dtos.NewErrorResponse(err))
		return
	}
	if errors.Is(err, &exceptions.BadRequest{}) {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// GetPredictedSerieById godoc
//
//	@Summary		Endpoint para obtener los valores de una serie pronosticadas por id
//	@Tags           Interfaz INA
//	@Produce		json
//	@Param          calibrado_id     path      int     true  "Id del calibrado"
//	@Param          serieId      query     int  true  "Id de la serie"  Format(string)
//	@Success		200	{object} dtos.CalibratedStreamsDataResponse
//	@Failure        400  {object}  dtos.ErrorResponse
//	@Failure        404  {object}  dtos.ErrorResponse
//	@Router			/series/pronosticadas/{calibrado_id} [get]
func (s inaControllerImpl) GetPredictedSerieById(ctx *gin.Context) {
	id, done := getUintPathParam(ctx, "calibrado_id")
	if done {
		return
	}
	streamId, done := getUintQueryParam(ctx, "serieId")
	if done {
		return
	}
	res, err := s.inaService.GetPredictedSerieById(id, streamId)
	if errors.Is(err, &exceptions.NotFound{}) {
		ctx.JSON(http.StatusNotFound, dtos.NewErrorResponse(err))
		return
	}
	if errors.Is(err, &exceptions.BadRequest{}) {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}
