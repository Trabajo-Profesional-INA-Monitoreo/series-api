package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
)

type WaterLevelsCalculator interface {
	AddMetrics(cards []dtos.MetricCard) []dtos.MetricCard
	Compute(float64)
	GetAlertsCount() uint64
	GetEvacuationCount() uint64
	GetLowWaterCount() uint64
}

type calculateWaterLevels struct {
	alertLevel           float64
	evacuationLevel      float64
	lowWaterLevel        float64
	countAlertLevel      float64
	countEvacuationLevel float64
	countLowWaterLevel   float64
}

type noWaterLevel struct {
}

const waterLevel = 2

func NewCalculatorOfWaterLevelsDependingOnVariable(station entities.Station, variableId uint64) WaterLevelsCalculator {
	if variableId != waterLevel {
		return &noWaterLevel{}
	}
	return &calculateWaterLevels{
		alertLevel:           station.AlertLevel,
		evacuationLevel:      station.EvacuationLevel,
		lowWaterLevel:        station.LowWaterLevel,
		countAlertLevel:      0,
		countEvacuationLevel: 0,
		countLowWaterLevel:   0,
	}
}

func NewCalculatorOfWaterLevels(alertLevel float64, evacuationLevel float64, lowWaterLevel float64) WaterLevelsCalculator {
	return &calculateWaterLevels{
		alertLevel:           alertLevel,
		evacuationLevel:      evacuationLevel,
		lowWaterLevel:        lowWaterLevel,
		countAlertLevel:      0,
		countEvacuationLevel: 0,
		countLowWaterLevel:   0,
	}
}

func (c *calculateWaterLevels) Compute(level float64) {
	if level >= c.evacuationLevel {
		c.countEvacuationLevel++
	} else if level >= c.alertLevel {
		c.countAlertLevel++
	} else if level <= c.lowWaterLevel {
		c.countLowWaterLevel++
	}
}

func (c *calculateWaterLevels) AddMetrics(metrics []dtos.MetricCard) []dtos.MetricCard {
	metrics = append(metrics, dtos.NewMetricCard(entities.MapMetricToString(entities.AguasBajas), c.countLowWaterLevel))
	metrics = append(metrics, dtos.NewMetricCard(entities.MapMetricToString(entities.AguasEvacuacion), c.countEvacuationLevel))
	return append(metrics, dtos.NewMetricCard(entities.MapMetricToString(entities.AguasAlerta), c.countAlertLevel))
}

func (n *calculateWaterLevels) GetAlertsCount() uint64 {
	return uint64(n.countAlertLevel)
}
func (n *calculateWaterLevels) GetEvacuationCount() uint64 {
	return uint64(n.countEvacuationLevel)
}
func (n *calculateWaterLevels) GetLowWaterCount() uint64 {
	return uint64(n.countLowWaterLevel)
}

func (n noWaterLevel) Compute(_ float64) {

}
func (n noWaterLevel) AddMetrics(cards []dtos.MetricCard) []dtos.MetricCard {
	return cards
}

func (n noWaterLevel) GetAlertsCount() uint64 {
	return 0
}
func (n noWaterLevel) GetEvacuationCount() uint64 {
	return 0
}
func (n noWaterLevel) GetLowWaterCount() uint64 {
	return 0
}
