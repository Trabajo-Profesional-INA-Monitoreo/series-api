package dtos

import "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"

type ConfiguredStream struct {
	ConfiguredStreamId      uint64              `gorm:"column:configured_stream_id" json:"ConfiguredStreamId"`
	StreamId                uint64              `binding:"required,gte=1" gorm:"column:stream_id" json:"StreamId"`
	StreamType              entities.StreamType `binding:"gte=0,lte=2" gorm:"column:stream_type" json:"StreamType"`
	UpdateFrequency         float64             `binding:"gte=0" gorm:"column:update_frequency" json:"UpdateFrequency"`
	CheckErrors             bool                `gorm:"column:check_errors" json:"CheckErrors"`
	UpperThreshold          float64             `gorm:"column:normal_upper_threshold" json:"UpperThreshold"`
	LowerThreshold          float64             `binding:"ltecsfield=UpperThreshold" gorm:"column:normal_lower_threshold" json:"LowerThreshold"`
	CalibrationId           uint64              `binding:"required_unless=StreamType 0" gorm:"column:calibration_id" json:"CalibrationId"`
	ForecastedRangeHours    *uint64             `binding:"excluded_unless=StreamType 1" json:"ForecastedRangeHours"`
	RedundanciesIds         *[]uint64           `gorm:"-" json:"RedundanciesIds"`
	Metrics                 *[]entities.Metric  `binding:"omitempty,gte=0,lte=5,unique,dive,gte=0,lte=4" gorm:"-" json:"Metrics"`
	ObservedRelatedStreamId *uint64             `binding:"excluded_if=StreamType 0" json:"ObservedRelatedStreamId"`
}

type ConfiguredStreamCreate struct {
	StreamId                uint64              `binding:"required,gte=1" gorm:"column:stream_id" json:"StreamId"`
	StreamType              entities.StreamType `binding:"lte=2" gorm:"column:stream_type" json:"StreamType"`
	UpdateFrequency         float64             `binding:"gte=0" gorm:"column:update_frequency" json:"UpdateFrequency"`
	CheckErrors             bool                `gorm:"column:check_errors" json:"CheckErrors"`
	UpperThreshold          float64             `gorm:"column:normal_upper_threshold" json:"UpperThreshold"`
	LowerThreshold          float64             `binding:"ltecsfield=UpperThreshold" gorm:"column:normal_lower_threshold" json:"LowerThreshold"`
	CalibrationId           uint64              `binding:"required_unless=StreamType 0" gorm:"column:calibration_id" json:"CalibrationId"`
	ForecastedRangeHours    *uint64             `binding:"excluded_unless=StreamType 1" json:"ForecastedRangeHours"`
	RedundanciesIds         []uint64            `gorm:"-" json:"RedundanciesIds"`
	Metrics                 []entities.Metric   `binding:"omitempty,gte=0,lte=5,unique,dive,gte=0,lte=4" gorm:"-" json:"Metrics"`
	ObservedRelatedStreamId *uint64             `binding:"excluded_if=StreamType 0" json:"ObservedRelatedStreamId"`
}
