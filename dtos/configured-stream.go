package dtos

import "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"

type ConfiguredStream struct {
	ConfiguredStreamId uint64             `gorm:"column:configured_stream_id"`
	StreamId           uint64             `gorm:"column:stream_id"`
	StreamType         uint64             `gorm:"column:stream_type"`
	UpdateFrequency    float64            `gorm:"column:update_frequency"`
	CheckErrors        bool               `gorm:"column:check_errors"`
	UpperThreshold     float64            `gorm:"column:normal_upper_threshold"`
	LowerThreshold     float64            `gorm:"column:normal_lower_threshold"`
	CalibrationId      uint64             `gorm:"column:calibration_id"`
	RedundanciesIds    *[]uint64          `gorm:"-"`
	Metrics            *[]entities.Metric `gorm:"-"`
}

type ConfiguredStreamCreate struct {
	StreamId        uint64            `gorm:"column:stream_id"`
	StreamType      uint64            `gorm:"column:stream_type"`
	UpdateFrequency float64           `gorm:"column:update_frequency"`
	CheckErrors     bool              `gorm:"column:check_errors"`
	UpperThreshold  float64           `gorm:"column:normal_upper_threshold"`
	LowerThreshold  float64           `gorm:"column:normal_lower_threshold"`
	CalibrationId   uint64            `gorm:"column:calibration_id"`
	RedundanciesIds []uint64          `gorm:"-"`
	Metrics         []entities.Metric `gorm:"-"`
}
