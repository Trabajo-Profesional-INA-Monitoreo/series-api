package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients/responses"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
)

func getMetricsForForecastedStream(data *responses.LastForecast, neededMetrics []entities.ConfiguredMetric) *[]dtos.MetricCard {
	return nil
}

func getMetricsForObservedOrCuratedStream(data []responses.ObservedDataResponse, neededMetrics []entities.ConfiguredMetric) *[]dtos.MetricCard {
	return nil
}
