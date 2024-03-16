package dtos

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"time"
)

type ErrorsCountPerDayAndType struct {
	ErrorTypeId entities.ErrorType `gorm:"column:error_type" json:"-"`
	ErrorType   string             `gorm:"-"`
	Total       int                `gorm:"column:total"`
	Date        time.Time          `gorm:"column:date"`
}

func (e *ErrorsCountPerDayAndType) ConvertToResponse() {
	e.ErrorType = entities.MapErrorTypeToString(e.ErrorTypeId)
}
