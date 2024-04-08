package controllers

import (
	"fmt"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func getDates(ctx *gin.Context) (time.Time, time.Time, bool) {
	timeStartQuery := ctx.DefaultQuery("timeStart", time.Now().Add(-DaysPerWeek*HoursPerDay*time.Hour).Format(time.DateOnly))
	timeEndQuery := ctx.DefaultQuery("timeEnd", time.Now().Add(HoursPerDay*time.Hour).Format(time.DateOnly))
	timeStart, err := time.Parse(time.DateOnly, timeStartQuery)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("error parsing time: %v", err)))
		return time.Time{}, time.Time{}, true
	}
	timeEnd, err := time.Parse(time.DateOnly, timeEndQuery)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("error parsing time: %v", err)))
		return time.Time{}, time.Time{}, true
	}
	return timeStart, timeEnd, false
}

func getUintQueryParam(ctx *gin.Context, queryParam string) (uint64, bool) {
	configurationIdQuery, sent := ctx.GetQuery(queryParam)
	if !sent {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf(queryParam+" missing")))
		return 0, true
	}
	configurationId, err := strconv.ParseUint(configurationIdQuery, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dtos.NewErrorResponse(fmt.Errorf("wrong format for "+queryParam+", should be uint")))
		return 0, true
	}
	return configurationId, false
}
