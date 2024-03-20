package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients/responses"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"sort"
)

func getMetricsForForecastedStream(data *responses.LastForecast, neededMetrics []entities.ConfiguredMetric) *[]dtos.MetricCard {
	if len(neededMetrics) == 0 {
		return nil
	}
	var metrics []dtos.MetricCard
	setUpFirstValues := false
	var minValue float64
	var maxValue float64
	sumOfValues := 0.0
	validValues := data.GetMainForecast()
	totalValues := len(validValues)
	for _, value := range validValues {
		if !setUpFirstValues {
			minValue = value
			maxValue = value
			setUpFirstValues = true
		}
		if minValue > value {
			minValue = value
		}
		if value > maxValue {
			maxValue = value
		}
		sumOfValues += value
		validValues = append(validValues, value)
	}

	for _, metric := range neededMetrics {
		metricName := entities.MapMetricToString(metric.MetricId)
		metricValue := 0.0
		if metric.MetricId == entities.Mediana {
			sort.Float64s(validValues)
			middle := totalValues / 2
			if len(validValues)%2 == 0 {
				metricValue = (validValues[middle-1] + validValues[middle]) / 2
			} else {
				metricValue = validValues[middle]
			}
		} else if metric.MetricId == entities.Maximo {
			metricValue = maxValue
		} else if metric.MetricId == entities.Minimo {
			metricValue = minValue
		} else if metric.MetricId == entities.Media {
			metricValue = sumOfValues / float64(totalValues)
		}
		metrics = append(metrics, dtos.NewMetricCard(metricName, metricValue))
	}

	metrics = append(metrics, dtos.NewMetricCard(entities.MapMetricToString(entities.Observaciones), float64(totalValues)))
	return &metrics
}

func getMetricsForObservedOrCuratedStream(data []responses.ObservedDataResponse, neededMetrics []entities.ConfiguredMetric) *[]dtos.MetricCard {
	if len(neededMetrics) == 0 {
		return nil
	}
	var metrics []dtos.MetricCard
	setUpFirstValues := false
	var minValue float64
	var maxValue float64
	nullValues := 0
	sumOfValues := 0.0
	totalValidValues := 0
	var validValues []float64
	for _, dataNode := range data {
		if dataNode.Value != nil {
			if !setUpFirstValues {
				minValue = *dataNode.Value
				maxValue = *dataNode.Value
				setUpFirstValues = true
			}
			if minValue > *dataNode.Value {
				minValue = *dataNode.Value
			}
			if *dataNode.Value > maxValue {
				maxValue = *dataNode.Value
			}
			sumOfValues += *dataNode.Value
			totalValidValues++
			validValues = append(validValues, *dataNode.Value)
		} else {
			nullValues++
		}
	}

	for _, metric := range neededMetrics {
		metricName := entities.MapMetricToString(metric.MetricId)
		metricValue := 0.0
		if metric.MetricId == entities.Mediana {
			sort.Float64s(validValues)
			middle := totalValidValues / 2
			if len(validValues)%2 == 0 {
				metricValue = (validValues[middle-1] + validValues[middle]) / 2
			} else {
				metricValue = validValues[middle]
			}
		} else if metric.MetricId == entities.Maximo {
			metricValue = maxValue
		} else if metric.MetricId == entities.Minimo {
			metricValue = minValue
		} else if metric.MetricId == entities.Media {
			metricValue = sumOfValues / float64(totalValidValues)
		} else if metric.MetricId == entities.Nulos {
			metricValue = float64(nullValues)
		}
		metrics = append(metrics, dtos.NewMetricCard(metricName, metricValue))
	}

	metrics = append(metrics, dtos.NewMetricCard(entities.MapMetricToString(entities.Observaciones), float64(totalValidValues)))
	return &metrics
}
