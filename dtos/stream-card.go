package dtos

type StreamCard struct {
	StreamId           uint64 `gorm:"stream_id"`
	ConfiguredStreamId uint64 `gorm:"configured_stream_id"`
	VariableId         uint64 `gorm:"variable_id"`
	VariableName       string `gorm:"variable_name"`
	ProcedureId        uint64 `gorm:"procedure_id"`
	ProcedureName      string `gorm:"procedure_name"`
	StationId          uint64 `gorm:"station_id"`
	StationName        string `gorm:"station_name"`
	CheckErrors        bool   `gorm:"check_errors"`
	TotalErrors        *uint64
}

type StreamCardsResponse struct {
	Content  *[]*StreamCard
	Pageable Pageable
}

func NewStreamCardsResponse(content []*StreamCard, pageable Pageable) *StreamCardsResponse {
	return &StreamCardsResponse{Content: &content, Pageable: pageable}
}
