package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients/responses"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type MockInaApiClient struct {
	ObservedData []responses.ObservedDataResponse
	Error        error
	LastForecast responses.LastForecast
}

func (m MockInaApiClient) GetLastForecast(calibrationId uint64) (*responses.LastForecast, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return &m.LastForecast, nil
}

func (m MockInaApiClient) GetObservedData(streamId uint64, timeStart time.Time, timeEnd time.Time) ([]responses.ObservedDataResponse, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.ObservedData, nil
}

func (m MockInaApiClient) GetStream(streamId uint64) (*responses.InaStreamResponse, error) {
	return nil, nil
}

func TestShouldReturnAllTheConfigurableMetricsOfAnObservedStream(t *testing.T) {
	neededMetrics := []entities.ConfiguredMetric{
		{MetricId: entities.Mediana},
		{MetricId: entities.Maximo},
		{MetricId: entities.Minimo},
		{MetricId: entities.Media},
		{MetricId: entities.Nulos},
	}
	var observedData []responses.ObservedDataResponse
	for i := 0; i < 10; i++ {
		level := float64(i)
		observedData = append(observedData, responses.ObservedDataResponse{Value: &level})
	}

	metrics := getMetricsForObservedOrCuratedStream(observedData, neededMetrics, &noWaterLevel{})
	assert.Equal(t, 6, len(*metrics))
}

func TestShouldReturnTheNullMetricOfAnObservedStream(t *testing.T) {
	neededMetrics := []entities.ConfiguredMetric{
		{MetricId: entities.Nulos},
	}
	var observedData []responses.ObservedDataResponse
	level := 1.0
	observedData = append(observedData, responses.ObservedDataResponse{Value: &level})
	observedData = append(observedData, responses.ObservedDataResponse{Value: nil})
	observedData = append(observedData, responses.ObservedDataResponse{Value: nil})

	metrics := getMetricsForObservedOrCuratedStream(observedData, neededMetrics, &noWaterLevel{})
	assert.Equal(t, 2, len(*metrics))
	assert.Equal(t, 2.0, (*metrics)[0].Value)
}

func TestShouldReturnTheMinValueMetricOfAnObservedStream(t *testing.T) {
	neededMetrics := []entities.ConfiguredMetric{
		{MetricId: entities.Minimo},
	}
	var observedData []responses.ObservedDataResponse
	for i := 5; i < 10; i++ {
		level := float64(i)
		observedData = append(observedData, responses.ObservedDataResponse{Value: &level})
	}
	for i := 3; i < 6; i++ {
		level := float64(i)
		observedData = append(observedData, responses.ObservedDataResponse{Value: &level})
	}

	metrics := getMetricsForObservedOrCuratedStream(observedData, neededMetrics, &noWaterLevel{})
	assert.Equal(t, 2, len(*metrics))
	assert.Equal(t, 3.0, (*metrics)[0].Value)
}

func TestShouldReturnTheMaxValueMetricOfAnObservedStream(t *testing.T) {
	neededMetrics := []entities.ConfiguredMetric{
		{MetricId: entities.Maximo},
	}
	var observedData []responses.ObservedDataResponse
	for i := 5; i < 10; i++ {
		level := float64(i)
		observedData = append(observedData, responses.ObservedDataResponse{Value: &level})
	}
	for i := 3; i < 6; i++ {
		level := float64(i)
		observedData = append(observedData, responses.ObservedDataResponse{Value: &level})
	}

	metrics := getMetricsForObservedOrCuratedStream(observedData, neededMetrics, &noWaterLevel{})
	assert.Equal(t, 2, len(*metrics))
	assert.Equal(t, 9.0, (*metrics)[0].Value)
}

func TestShouldReturnAllTheConfigurableMetricsOfAForecastedStream(t *testing.T) {
	neededMetrics := []entities.ConfiguredMetric{
		{MetricId: entities.Mediana},
		{MetricId: entities.Maximo},
		{MetricId: entities.Minimo},
		{MetricId: entities.Media},
	}
	lastForecast := &responses.Forecast{
		MainForecast: &responses.ForecastedStream{Forecasts: [][]string{{"", "", "1.0", ""}, {"", "", "1.1", ""}, {"", "", "1.2", ""}, {"", "", "1.5", ""}}},
	}

	metrics := getMetricsForForecastedStream(lastForecast, neededMetrics, &noWaterLevel{})
	assert.Equal(t, 5, len(*metrics))
}

func TestShouldReturnMinValueOfAForecastedStream(t *testing.T) {
	neededMetrics := []entities.ConfiguredMetric{
		{MetricId: entities.Minimo},
	}
	lastForecast := &responses.Forecast{
		MainForecast: &responses.ForecastedStream{Forecasts: [][]string{{"", "", "1.0", ""}, {"", "", "0.5", ""}, {"", "", "1.2", ""}, {"", "", "1.5", ""}}},
	}

	metrics := getMetricsForForecastedStream(lastForecast, neededMetrics, &noWaterLevel{})
	assert.Equal(t, 2, len(*metrics))
	assert.Equal(t, 0.5, (*metrics)[0].Value)
}

func TestShouldReturnMaxValueOfAForecastedStream(t *testing.T) {
	neededMetrics := []entities.ConfiguredMetric{
		{MetricId: entities.Maximo},
	}
	lastForecast := &responses.Forecast{
		MainForecast: &responses.ForecastedStream{Forecasts: [][]string{{"", "", "1.0", ""}, {"", "", "0.5", ""}, {"", "", "1.2", ""}, {"", "", "1.5", ""}}},
	}

	metrics := getMetricsForForecastedStream(lastForecast, neededMetrics, &noWaterLevel{})
	assert.Equal(t, 2, len(*metrics))
	assert.Equal(t, 1.5, (*metrics)[0].Value)
}

func TestShouldReturnMedianValueOfAForecastedStream(t *testing.T) {
	neededMetrics := []entities.ConfiguredMetric{
		{MetricId: entities.Mediana},
	}
	lastForecast := &responses.Forecast{
		MainForecast: &responses.ForecastedStream{Forecasts: [][]string{{"", "", "1.0", ""}, {"", "", "2.0", ""}, {"", "", "3.0", ""}, {"", "", "4.0", ""}}},
	}

	metrics := getMetricsForForecastedStream(lastForecast, neededMetrics, &noWaterLevel{})
	assert.Equal(t, 2, len(*metrics))
	assert.Equal(t, 2.5, (*metrics)[0].Value)
}

func TestShouldReturnAvgValueOfAForecastedStream(t *testing.T) {
	neededMetrics := []entities.ConfiguredMetric{
		{MetricId: entities.Media},
	}
	lastForecast := &responses.Forecast{
		MainForecast: &responses.ForecastedStream{Forecasts: [][]string{{"", "", "1.0", ""}, {"", "", "2.0", ""}, {"", "", "3.0", ""}}},
	}

	metrics := getMetricsForForecastedStream(lastForecast, neededMetrics, &noWaterLevel{})
	assert.Equal(t, 2, len(*metrics))
	assert.Equal(t, float64(2), (*metrics)[0].Value)
}
