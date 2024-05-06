package dtos

type BehaviourStream struct {
	StreamId        uint64   `json:"StreamId"`
	AlertLevel      *float64 `json:"AlertLevel"`
	EvacuationLevel *float64 `json:"EvacuationLevel"`
	LowWaterLevel   *float64 `json:"LowWaterLevel"`
}

type BehaviourStreamsResponse struct {
	TotalValuesCount     uint64 `json:"TotalValuesCount"`
	CountAlertLevel      uint64 `json:"CountAlertLevel"`
	CountEvacuationLevel uint64 `json:"CountEvacuationLevel"`
	CountLowWaterLevel   uint64 `json:"CountLowWaterLevel"`
}

func NewBehaviourStreamsResponse() *BehaviourStreamsResponse {
	return &BehaviourStreamsResponse{
		TotalValuesCount:     0,
		CountAlertLevel:      0,
		CountEvacuationLevel: 0,
		CountLowWaterLevel:   0,
	}
}
