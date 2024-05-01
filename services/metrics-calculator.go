package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients/responses"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	log "github.com/sirupsen/logrus"
	"sort"
	"time"
)

type metricParameters struct {
	setUpFirstValues bool
	validValues      []float64
	maxValue         float64
	minValue         float64
	sumOfValues      float64
	nullValues       int
}

func getMetricsForForecastedStream(data *responses.Forecast, neededMetrics []entities.ConfiguredMetric, waterLevelCalculator WaterLevelsCalculator) *[]dtos.MetricCard {
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

func getMetricsForObservedOrCuratedStream(data []responses.ObservedDataResponse, neededMetrics []entities.ConfiguredMetric, waterLevelCalculator WaterLevelsCalculator) *[]dtos.MetricCard {
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
			metricValue = metricsValues.sumOfValues / float64(totalValidValues)
		} else if metric.MetricId == entities.Nulos {
			metricValue = float64(metricsValues.nullValues)
		}
		metrics = append(metrics, dtos.NewMetricCard(metricName, metricValue))
	}
	metrics = waterLevelCalculator.AddMetrics(metrics)
	metrics = append(metrics, dtos.NewMetricCard(entities.MapMetricToString(entities.Observaciones), float64(totalValidValues)))
	return &metrics
}

func calculateMedian(metricsValues *metricParameters, totalValidValues int) float64 {
	sort.Float64s(metricsValues.validValues)
	middle := totalValidValues / 2
	if totalValidValues%2 == 0 {
		return (metricsValues.validValues[middle-1] + metricsValues.validValues[middle]) / 2
	}
	return metricsValues.validValues[middle]
}

func getLevelsCountForAllStreams(behaviourStreams []dtos.BehaviourStream, timeStart time.Time, timeEnd time.Time, inaApiClient clients.InaAPiClient) *dtos.BehaviourStreamsResponse {
	behaviourResponse := dtos.NewBehaviourStreamsResponse()
	for _, stream := range behaviourStreams {
		values, err := inaApiClient.GetObservedData(stream.StreamId, timeStart, timeEnd)
		if err != nil {
			log.Errorf("GetOutputBehaviourMetrics | Could not get metrics for stream with id %v: %v", stream.StreamId, err)
			continue
		}
		calculator := NewCalculatorOfWaterLevels(stream.AlertLevel, stream.EvacuationLevel, stream.LowWaterLevel)
		for _, observedData := range values {
			if observedData.Value != nil {
				calculator.Compute(*observedData.Value)
				behaviourResponse.TotalValuesCount += 1
			}
		}
		behaviourResponse.CountAlertLevel += calculator.GetAlertsCount()
		behaviourResponse.CountLowWaterLevel += calculator.GetLowWaterCount()
		behaviourResponse.CountEvacuationLevel += calculator.GetEvacuationCount()
	}
	return behaviourResponse
}
