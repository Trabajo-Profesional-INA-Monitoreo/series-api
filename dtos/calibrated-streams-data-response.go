package dtos

import "time"

type CalibratedStreamsDataResponse struct {
	P05Streams  []CalibratedStreamsData
	MainStreams []CalibratedStreamsData
	P75Streams  []CalibratedStreamsData
	P95Streams  []CalibratedStreamsData
	P25Streams  []CalibratedStreamsData
}

type CalibratedStreamsData struct {
	Time      time.Time
	Value     float64
	Qualifier string
}
