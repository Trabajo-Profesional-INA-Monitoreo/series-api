package metrics_service

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
	GetStreamLevels(id uint64) []dtos.StreamLevel
}

type calculateWaterLevels struct {
	alertLevel           *float64
	evacuationLevel      *float64
	lowWaterLevel        *float64
	countAlertLevel      float64
	countEvacuationLevel float64
	countLowWaterLevel   float64
}

type noWaterLevel struct {
}

const waterLevel = 2
const alertLevel = "Alerta"
const evacuationLevel = "Evacuación"
const lowWaterLevel = "Aguas Bajas"

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

func NewCalculatorOfWaterLevels(alertLevel *float64, evacuationLevel *float64, lowWaterLevel *float64) WaterLevelsCalculator {
	return &calculateWaterLevels{
		alertLevel:           alertLevel,
		evacuationLevel:      evacuationLevel,
		lowWaterLevel:        lowWaterLevel,
		countAlertLevel:      0,
		countEvacuationLevel: 0,
		countLowWaterLevel:   0,
	}
}

func (c *calculateWaterLevels) GetStreamLevels(id uint64) []dtos.StreamLevel {
	var levels []dtos.StreamLevel
	if c.countEvacuationLevel > 0 {
		levels = append(levels, dtos.StreamLevel{
			StreamId: id,
			Level:    evacuationLevel,
		})
	}
	if c.countAlertLevel > 0 {
		levels = append(levels, dtos.StreamLevel{
			StreamId: id,
			Level:    alertLevel,
		})
	}
	if c.countLowWaterLevel > 0 {
		levels = append(levels, dtos.StreamLevel{
			StreamId: id,
			Level:    lowWaterLevel,
		})
	}
	return levels
}

func (c *calculateWaterLevels) Compute(level float64) {
	if c.evacuationLevel != nil && level >= *c.evacuationLevel {
		c.countEvacuationLevel++
	} else if c.alertLevel != nil && level >= *c.alertLevel {
		c.countAlertLevel++
	} else if c.lowWaterLevel != nil && level <= *c.lowWaterLevel {
		c.countLowWaterLevel++
	}
}

func (c *calculateWaterLevels) AddMetrics(metrics []dtos.MetricCard) []dtos.MetricCard {
	if c.lowWaterLevel != nil {
		metrics = append(metrics, dtos.NewMetricCard(entities.MapMetricToString(entities.AguasBajas), c.countLowWaterLevel))
	}
	if c.alertLevel != nil {
		metrics = append(metrics, dtos.NewMetricCard(entities.MapMetricToString(entities.AguasAlerta), c.countAlertLevel))
	}
	if c.evacuationLevel != nil {
		metrics = append(metrics, dtos.NewMetricCard(entities.MapMetricToString(entities.AguasEvacuacion), c.countEvacuationLevel))
	}
	return metrics
}

func (c *calculateWaterLevels) GetAlertsCount() uint64 {
	return uint64(c.countAlertLevel)
}
func (c *calculateWaterLevels) GetEvacuationCount() uint64 {
	return uint64(c.countEvacuationLevel)
}
func (c *calculateWaterLevels) GetLowWaterCount() uint64 {
	return uint64(c.countLowWaterLevel)
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

func (n noWaterLevel) GetStreamLevels(id uint64) []dtos.StreamLevel {
	return nil
}
