package dtos

import "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"

type StreamCard struct {
	StreamId                uint64              `gorm:"stream_id" json:"StreamId"`
	ConfiguredStreamId      uint64              `gorm:"configured_stream_id" json:"ConfiguredStreamId"`
	VariableId              uint64              `gorm:"variable_id" json:"VariableId"`
	VariableName            string              `gorm:"variable_name" json:"VariableName"`
	ProcedureId             uint64              `gorm:"procedure_id" json:"ProcedureId"`
	ProcedureName           string              `gorm:"procedure_name" json:"ProcedureName"`
	StationId               uint64              `gorm:"station_id" json:"StationId"`
	StationName             string              `gorm:"station_name" json:"StationName"`
	CheckErrors             bool                `gorm:"check_errors" json:"CheckErrors"`
	TotalErrors             *uint64             `json:"TotalErrors"`
	StreamType              entities.StreamType `json:"StreamType"`
	CalibrationId           uint64              `json:"CalibrationId"`
	ObservedRelatedStreamId *uint64             `json:"ObservedRelatedStreamId"`
}

type StreamCardsResponse struct {
	Content  *[]*StreamCard `json:"Content"`
	Pageable Pageable       `json:"Pageable"`
}

func NewStreamCardsResponse(content []*StreamCard, pageable Pageable) *StreamCardsResponse {
	return &StreamCardsResponse{Content: &content, Pageable: pageable}
}
