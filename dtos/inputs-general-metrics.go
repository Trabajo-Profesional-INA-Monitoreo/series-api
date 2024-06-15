package dtos

type InputsGeneralMetrics struct {
	TotalStreams  int `json:"TotalStreams"`
	TotalStations int `json:"TotalStations"`
}

type TotalStreamsWithNullValues struct {
	TotalStreams         int     `json:"TotalStreams"`
	TotalStreamsWithNull int     `json:"TotalStreamsWithNull"`
	Streams              []int64 `json:"Streams"`
}

type TotalStreamsWithObservedOutlier struct {
	TotalStreams                    int     `json:"TotalStreams"`
	TotalStreamsWithObservedOutlier int     `json:"TotalStreamsWithObservedOutlier"`
	Streams                         []int64 `json:"Streams"`
}

type TotalStreamsWithDelay struct {
	TotalStreams          int     `json:"TotalStreams"`
	TotalStreamsWithDelay int     `json:"TotalStreamsWithDelay"`
	Streams               []int64 `json:"Streams"`
}
