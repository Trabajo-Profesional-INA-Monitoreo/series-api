package dtos

type BehaviourStream struct {
	StreamId        uint64
	AlertLevel      float64
	EvacuationLevel float64
	LowWaterLevel   float64
}

type BehaviourStreamsResponse struct {
	TotalValuesCount     uint64
	CountAlertLevel      uint64
	CountEvacuationLevel uint64
	CountLowWaterLevel   uint64
}

func NewBehaviourStreamsResponse() *BehaviourStreamsResponse {
	return &BehaviourStreamsResponse{
		TotalValuesCount:     0,
		CountAlertLevel:      0,
		CountEvacuationLevel: 0,
		CountLowWaterLevel:   0,
	}
}
