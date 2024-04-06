package dtos

type ErrorsPerConfigStream struct {
	ConfiguredStreamId uint64 `gorm:"configured_stream_id"`
	ErrorsCount        uint64 `gorm:"errors_count"`
}
