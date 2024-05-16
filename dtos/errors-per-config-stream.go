package dtos

type ErrorsPerConfigStream struct {
	ConfiguredStreamId uint64 `gorm:"configured_stream_id" json:"ConfiguredStreamId"`
	ErrorsCount        uint64 `gorm:"errors_count" json:"ErrorsCount"`
}
