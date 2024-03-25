package dtos

type StreamCard struct {
	StreamId           uint64
	ConfiguredStreamId uint64
	VarId              string
	VariableName       string
	ProcId             uint64
	ProcedureName      string
	StationId          uint64
	StationName        string
	CheckErrors        bool
	TotalErrors        *uint64
}

type StreamCardsResponse struct {
	Content  []StreamCard
	Pageable Pageable
}

func NewStreamCardsResponse(content []StreamCard, pageable Pageable) *StreamCardsResponse {
	return &StreamCardsResponse{Content: content, Pageable: pageable}
}
