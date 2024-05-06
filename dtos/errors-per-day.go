package dtos

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"time"
)

type ErrorsCountPerDayAndType struct {
	ErrorTypeId entities.ErrorType `gorm:"column:error_type" json:"-"`
	ErrorType   string             `gorm:"-" json:"ErrorType"`
	Total       int                `gorm:"column:total" json:"Total"`
	Date        time.Time          `gorm:"column:date" json:"Date"`
}

func (e *ErrorsCountPerDayAndType) ConvertToResponse() {
	e.ErrorType = entities.MapErrorTypeToString(e.ErrorTypeId)
}
