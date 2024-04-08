package dtos

type InputsGeneralMetrics struct {
	TotalStreams  int
	TotalStations int
	TotalNetworks int
}

type TotalStreamsWithNullValues struct {
	TotalStreams         int
	TotalStreamsWithNull int
}
