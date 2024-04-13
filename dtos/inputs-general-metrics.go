package dtos

type InputsGeneralMetrics struct {
	TotalStreams  int
	TotalStations int
}

type TotalStreamsWithNullValues struct {
	TotalStreams         int
	TotalStreamsWithNull int
}
