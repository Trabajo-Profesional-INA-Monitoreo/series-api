package responses

import "time"

type MainStreamFromStation struct {
	StreamId   uint64
	LastUpdate *time.Time
	AlertLevel *float64
}

func NewMainStreamFromStation(response InaStreamResponse) *MainStreamFromStation {
	return &MainStreamFromStation{uint64(response.Id), response.DateRange.TimeEnd, response.Station.AlertLevel}
}
