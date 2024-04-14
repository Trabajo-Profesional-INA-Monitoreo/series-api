package dtos

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"time"
)

type DetectedErrorsOfStream struct {
	Content  []*ErrorDto
	Pageable Pageable
}

func NewDetectedErrorsOfStream(content []*ErrorDto, pageable Pageable) *DetectedErrorsOfStream {
	return &DetectedErrorsOfStream{Content: content, Pageable: pageable}
}

type ErrorDto struct {
	ErrorId       uint64
	DetectedDate  time.Time
	ErrorTypeId   entities.ErrorType
	ErrorTypeName string
	ExtraInfo     string
}

func (d *DetectedErrorsOfStream) ConvertToResponse() {
	for _, dto := range d.Content {
		dto.ErrorTypeName = entities.MapErrorTypeToString(dto.ErrorTypeId)
	}

}
