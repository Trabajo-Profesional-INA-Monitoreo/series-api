package dtos

type InputsGeneralMetrics struct {
	TotalStreams  int `json:"TotalStreams"`
	TotalStations int `json:"TotalStations"`
}

type TotalStreamsWithNullValues struct {
	TotalStreams         int `json:"TotalStreams"`
	TotalStreamsWithNull int `json:"TotalStreamsWithNull"`
}

type TotalStreamsWithObservedOutlier struct {
	TotalStreams                    int `json:"TotalStreams"`
	TotalStreamsWithObservedOutlier int `json:"TotalStreamsWithObservedOutlier"`
}
