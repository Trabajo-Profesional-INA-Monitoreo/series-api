package dtos

type StreamsPerStationResponse struct {
	Stations []*StreamsPerStation
}

type StreamsPerStation struct {
	StationName  string `gorm:"column:station_name" json:"StationName"`
	StationId    string `gorm:"column:station_id" json:"StationId"`
	StreamsCount int    `gorm:"column:streams_count" json:"StreamsCount"`
	ErrorCount   int    `json:"ErrorCount"`
}

type ErrorsOfStations struct {
	StationId  string `gorm:"column:station_id" json:"StationId"`
	ErrorCount int    `gorm:"column:error_count" json:"ErrorCount"`
}
