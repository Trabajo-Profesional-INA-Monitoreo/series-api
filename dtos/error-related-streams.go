package dtos

type ErrorRelatedStream struct {
	StreamId    uint64 `json:"StreamId"`
	StationName string `json:"StationName"`
	StationId   uint64 `json:"StationId"`
}
