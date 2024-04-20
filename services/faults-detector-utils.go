package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients/responses"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"strconv"
	"time"
)

func shouldFetchObservedStream(forecast *responses.Forecast, configuredStream entities.ConfiguredStream) bool {
	return (forecast.P05Forecast != nil || forecast.P95Forecast != nil) && configuredStream.ObservedRelatedStreamId != nil
}

func valueIsAnOutlier(configuredStream entities.ConfiguredStream, observed responses.ObservedDataResponse) bool {
	return configuredStream.NormalLowerThreshold > *observed.Value || configuredStream.NormalUpperThreshold < *observed.Value
}

func contains(configuredStreams []entities.ConfiguredStream, configuredStream entities.ConfiguredStream) bool {
	for _, cs := range configuredStreams {
		if cs.ConfiguredStreamId == configuredStream.ConfiguredStreamId {
			return true
		}
	}
	return false
}

func getDateRangeOfForecast(forecast *responses.ForecastedStream) (time.Time, time.Time) {
	values := forecast.Forecasts
	timeString := values[0][0]
	timeStart, _ := time.Parse("2006-01-02T15:04:05Z07:00", timeString)
	timeString = values[len(values)-1][0]
	timeEnd, _ := time.Parse("2006-01-02T15:04:05Z07:00", timeString)
	return timeStart, timeEnd
}

func parseForecastedDate(forecast string) time.Time {
	timeForecast, _ := time.Parse("2006-01-02T15:04:05Z07:00", forecast)
	return timeForecast
}

func observedIsOutsideErrorBands(forecast *responses.Forecast, forecastedIndex int, observed *float64) bool {
	if observed == nil {
		return false
	}
	if forecast.P05Forecast != nil {
		value, _ := strconv.ParseFloat(forecast.P05Forecast.Forecasts[forecastedIndex][2], 64)
		if *observed < value {
			return true
		}
	}
	if forecast.P95Forecast != nil {
		value, _ := strconv.ParseFloat(forecast.P95Forecast.Forecasts[forecastedIndex][2], 64)
		if *observed > value {
			return true
		}
	}
	return false
}

func tooManyValuesOutsideBands(outsideBands int, totalObservedValues int) bool {
	return float64(outsideBands) > float64(totalObservedValues)*0.7
}
