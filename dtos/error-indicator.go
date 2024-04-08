package dtos

import "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"

type ErrorIndicator struct {
	ErrorTypeId entities.ErrorType `gorm:"column:error_type" json:"ErrorId"`
	ErrorType   string             `gorm:"-"`
	Count       int
}

func (e *ErrorIndicator) ConvertToResponse() {
	e.ErrorType = entities.MapErrorTypeToString(e.ErrorTypeId)
}
