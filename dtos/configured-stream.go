package dtos

import "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"

type ConfiguredStream struct {
	ConfiguredStreamId      uint64              `gorm:"column:configured_stream_id"`
	StreamId                uint64              `binding:"required,gte=1" gorm:"column:stream_id"`
	StreamType              entities.StreamType `binding:"gte=0,lte=2" gorm:"column:stream_type"`
	UpdateFrequency         float64             `binding:"gte=0" gorm:"column:update_frequency"`
	CheckErrors             bool                `gorm:"column:check_errors"`
	UpperThreshold          float64             `gorm:"column:normal_upper_threshold"`
	LowerThreshold          float64             `binding:"ltecsfield=UpperThreshold" gorm:"column:normal_lower_threshold"`
	CalibrationId           uint64              `binding:"excluded_if=StreamType 0" gorm:"column:calibration_id"`
	RedundanciesIds         *[]uint64           `gorm:"-"`
	Metrics                 *[]entities.Metric  `binding:"omitempty,gte=0,lte=4,dive,gte=0,lte=4" gorm:"-"`
	ObservedRelatedStreamId *uint64             `binding:"excluded_unless=StreamType 1"`
}

type ConfiguredStreamCreate struct {
	StreamId                uint64              `binding:"required,gte=1" gorm:"column:stream_id"`
	StreamType              entities.StreamType `binding:"lte=2" gorm:"column:stream_type"`
	UpdateFrequency         float64             `binding:"gte=0" gorm:"column:update_frequency"`
	CheckErrors             bool                `gorm:"column:check_errors"`
	UpperThreshold          float64             `gorm:"column:normal_upper_threshold"`
	LowerThreshold          float64             `binding:"ltecsfield=UpperThreshold" gorm:"column:normal_lower_threshold"`
	CalibrationId           uint64              `binding:"excluded_if=StreamType 0" gorm:"column:calibration_id"`
	RedundanciesIds         []uint64            `gorm:"-"`
	Metrics                 []entities.Metric   `binding:"omitempty,gte=0,lte=4,dive,gte=0,lte=4" gorm:"-"`
	ObservedRelatedStreamId *uint64             `binding:"excluded_unless=StreamType 1"`
}
