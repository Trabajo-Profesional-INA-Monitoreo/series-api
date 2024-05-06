package dtos

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"time"
)

type DetectedErrorsOfStream struct {
	Content  []*ErrorDto `json:"Content"`
	Pageable Pageable    `json:"Pageable"`
}

func NewDetectedErrorsOfStream(content []*ErrorDto, pageable Pageable) *DetectedErrorsOfStream {
	return &DetectedErrorsOfStream{Content: content, Pageable: pageable}
}

type ErrorDto struct {
	ErrorId       uint64             `json:"ErrorId"`
	DetectedDate  time.Time          `json:"DetectedDate"`
	ErrorTypeId   entities.ErrorType `json:"ErrorTypeId"`
	ErrorTypeName string             `json:"ErrorTypeName"`
	ExtraInfo     string             `json:"ExtraInfo"`
}

func (d *DetectedErrorsOfStream) ConvertToResponse() {
	for _, dto := range d.Content {
		dto.ErrorTypeName = entities.MapErrorTypeToString(dto.ErrorTypeId)
	}

}
