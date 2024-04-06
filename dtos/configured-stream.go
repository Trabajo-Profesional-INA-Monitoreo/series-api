package dtos

import "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"

type ConfiguredStream struct {
	StreamId        uint64
	StreamType      uint64
	UpdateFrequency float64
	CheckErrors     bool
	UpperThreshold  uint64
	LowerThreshold  uint64
	CalibrationId   uint64
	RedundanciesIds []uint64
	Metrics         []entities.Metric
}
