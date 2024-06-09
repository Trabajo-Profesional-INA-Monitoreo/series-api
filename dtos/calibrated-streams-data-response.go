package dtos

import "time"

type CalibratedStreamsDataResponse struct {
	P05Streams  []CalibratedStreamsData `json:"P05Streams"`
	MainStreams []CalibratedStreamsData `json:"MainStreams"`
	P75Streams  []CalibratedStreamsData `json:"P75Streams"`
	P95Streams  []CalibratedStreamsData `json:"P95Streams"`
	P25Streams  []CalibratedStreamsData `json:"P25Streams"`
	P99Streams  []CalibratedStreamsData `json:"P99Streams"`
	P01Streams  []CalibratedStreamsData `json:"P01Streams"`
}

type CalibratedStreamsData struct {
	Time  time.Time `json:"Time"`
	Value float64   `json:"Value"`
}
