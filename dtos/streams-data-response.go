package dtos

import "time"

type StreamsDataResponse struct {
	Streams []StreamsData
}

type StreamsData struct {
	Time  time.Time
	Value *float64
}
