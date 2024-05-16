package metrics_service

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients/responses"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"sort"
)

type metricParameters struct {
	setUpFirstValues bool
	validValues      []float64
	maxValue         float64
	minValue         float64
	sumOfValues      float64
	nullValues       int
}

func GetMetricsForForecastedStream(data *responses.Forecast, neededMetrics []entities.ConfiguredMetric, waterLevelCalculator WaterLevelsCalculator) *[]dtos.MetricCard {
	metricsValues := &metricParameters{
		setUpFirstValues: false,
		validValues:      responses.ConvertToFloats(data.MainForecast.Forecasts),
		maxValue:         0,
		minValue:         0,
		sumOfValues:      0,
	}

	for _, waterLevel := range metricsValues.validValues {
		calculateMetrics(metricsValues, waterLevel, waterLevelCalculator)
	}

	return addMetricsCards(neededMetrics, metricsValues, waterLevelCalculator)
}

func GetMetricsForObservedOrCuratedStream(data []responses.ObservedDataResponse, neededMetrics []entities.ConfiguredMetric, waterLevelCalculator WaterLevelsCalculator) *[]dtos.MetricCard {
	metricsValues := &metricParameters{
		setUpFirstValues: false,
		validValues:      []float64{},
		maxValue:         0,
		minValue:         0,
		sumOfValues:      0,
		nullValues:       0,
	}
	for _, dataNode := range data {
		if dataNode.Value == nil {
			metricsValues.nullValues++
			continue
		}
		calculateMetrics(metricsValues, *dataNode.Value, waterLevelCalculator)
		metricsValues.validValues = append(metricsValues.validValues, *dataNode.Value)
	}

	return addMetricsCards(neededMetrics, metricsValues, waterLevelCalculator)
}

func calculateMetrics(metricsValues *metricParameters, waterLevel float64, waterLevelCalculator WaterLevelsCalculator) {
	if !metricsValues.setUpFirstValues {
		metricsValues.minValue = waterLevel
		metricsValues.maxValue = waterLevel
		metricsValues.setUpFirstValues = true
	}
	if metricsValues.minValue > waterLevel {
		metricsValues.minValue = waterLevel
	}
	if waterLevel > metricsValues.maxValue {
		metricsValues.maxValue = waterLevel
	}
	waterLevelCalculator.Compute(waterLevel)
	metricsValues.sumOfValues += waterLevel
}

func addMetricsCards(neededMetrics []entities.ConfiguredMetric, metricsValues *metricParameters, waterLevelCalculator WaterLevelsCalculator) *[]dtos.MetricCard {
	var metrics []dtos.MetricCard
	totalValidValues := len(metricsValues.validValues)
	for _, metric := range neededMetrics {
		metricName := entities.MapMetricToString(metric.MetricId)
		metricValue := 0.0
		if metric.MetricId == entities.Mediana {
			metricValue = calculateMedian(metricsValues, totalValidValues)
		} else if metric.MetricId == entities.Maximo {
			metricValue = metricsValues.maxValue
		} else if metric.MetricId == entities.Minimo {
			metricValue = metricsValues.minValue
		} else if metric.MetricId == entities.Media {
			metricValue = calculateAverage(totalValidValues, metricsValues)
		} else if metric.MetricId == entities.Nulos {
			metricValue = float64(metricsValues.nullValues)
		}
		metrics = append(metrics, dtos.NewMetricCard(metricName, metricValue))
	}
	metrics = waterLevelCalculator.AddMetrics(metrics)
	metrics = append(metrics, dtos.NewMetricCard(entities.MapMetricToString(entities.Observaciones), float64(totalValidValues)))
	return &metrics
}

func calculateAverage(totalValidValues int, metricsValues *metricParameters) float64 {
	if totalValidValues == 0 {
		return 0
	}
	return metricsValues.sumOfValues / float64(totalValidValues)
}

func calculateMedian(metricsValues *metricParameters, totalValidValues int) float64 {
	if totalValidValues == 0 {
		return 0
	}
	sort.Float64s(metricsValues.validValues)
	middle := totalValidValues / 2
	if middle == 0 {
		return metricsValues.validValues[0]
	}
	if totalValidValues%2 == 0 {
		return (metricsValues.validValues[middle-1] + metricsValues.validValues[middle]) / 2
	}
	return metricsValues.validValues[middle]
}
