package dtos

type StreamsPerStationResponse struct {
	Stations []StreamsPerStation
}

type StreamsPerStation struct {
	StationName  string `gorm:"column:stationname"`
	StationId    string `gorm:"column:stationid"`
	StreamsCount int    `gorm:"column:streamscount"`
}
