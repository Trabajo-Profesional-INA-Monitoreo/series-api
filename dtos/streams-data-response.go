package dtos

import "time"

type StreamsDataResponse struct {
	Streams []StreamsData `json:"Streams"`
}

type StreamsData struct {
	Time  time.Time `json:"Time"`
	Value *float64  `json:"Value"`
}
