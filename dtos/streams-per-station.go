package dtos

import "time"

type StreamsPerStationResponse struct {
	Stations []*StreamsPerStation `json:"Stations"`
	Pageable Pageable             `json:"Pageable"`
}

type StreamsPerStation struct {
	StationName      string     `gorm:"column:station_name" json:"StationName"`
	StationId        string     `gorm:"column:station_id" json:"StationId"`
	StreamsCount     int        `gorm:"column:streams_count" json:"StreamsCount"`
	ErrorCount       int        `json:"ErrorCount"`
	MainStreamId     *uint64    `json:"MainStreamId"`
	AlertWaterLevels uint64     `json:"AlertWaterLevels"`
	TotalWaterLevels uint64     `json:"TotalWaterLevels"`
	LastUpdate       *time.Time `json:"LastUpdate"`
}

type ErrorsOfStations struct {
	StationId  string `gorm:"column:station_id" json:"StationId"`
	ErrorCount int    `gorm:"column:error_count" json:"ErrorCount"`
}
