package metrics_service

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShouldntAddMetricsIfLEvelsAreNull(t *testing.T) {
	calc := NewCalculatorOfWaterLevels(nil, nil, nil)
	var cards []dtos.MetricCard

	result := calc.AddMetrics(cards)

	assert.Empty(t, result)
}

func TestShouldAddALowWaterLevelValueToTheCount(t *testing.T) {
	lowLevel := 1.0
	calc := NewCalculatorOfWaterLevels(nil, nil, &lowLevel)
	calc.Compute(0.5)
	assert.Equal(t, uint64(1), calc.GetLowWaterCount())
	assert.Equal(t, uint64(0), calc.GetAlertsCount())
	assert.Equal(t, uint64(0), calc.GetEvacuationCount())
}

func TestShouldAddAlertWaterLevelValueToTheCount(t *testing.T) {
	alertLevel := 1.0
	calc := NewCalculatorOfWaterLevels(&alertLevel, nil, nil)
	calc.Compute(1.5)
	assert.Equal(t, uint64(0), calc.GetLowWaterCount())
	assert.Equal(t, uint64(1), calc.GetAlertsCount())
	assert.Equal(t, uint64(0), calc.GetEvacuationCount())
}

func TestShouldAddEvacuationWaterLevelValueToTheCount(t *testing.T) {
	evacLevel := 1.0
	calc := NewCalculatorOfWaterLevels(nil, &evacLevel, nil)
	calc.Compute(1.5)
	assert.Equal(t, uint64(0), calc.GetLowWaterCount())
	assert.Equal(t, uint64(0), calc.GetAlertsCount())
	assert.Equal(t, uint64(1), calc.GetEvacuationCount())
}

func TestShouldAddTheMetricsToTheCardsList(t *testing.T) {
	evacLevel := 2.0
	alertLevel := 1.0
	lowLevel := 0.5
	calc := NewCalculatorOfWaterLevels(&alertLevel, &evacLevel, &lowLevel)
	var cards []dtos.MetricCard

	calc.Compute(2.1)
	calc.Compute(2.2)
	calc.Compute(1.9)
	calc.Compute(1.5)
	calc.Compute(0.5)

	result := calc.AddMetrics(cards)

	assert.Equal(t, 3, len(result))
	assert.Equal(t, uint64(2), calc.GetAlertsCount())
	assert.Equal(t, uint64(2), calc.GetEvacuationCount())
	assert.Equal(t, uint64(1), calc.GetLowWaterCount())
}
