package dtos

type StreamsPerStationResponse struct {
	Stations []*StreamsPerStation
}

type StreamsPerStation struct {
	StationName  string `gorm:"column:station_name"`
	StationId    string `gorm:"column:station_id"`
	StreamsCount int    `gorm:"column:streams_count"`
	ErrorCount   int
}

type ErrorsOfStations struct {
	StationId  string `gorm:"column:station_id"`
	ErrorCount int    `gorm:"column:error_count"`
}
