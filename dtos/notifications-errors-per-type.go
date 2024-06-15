package dtos

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
)

type NotificationsErrorsCountPerType struct {
	ErrorTypeId     entities.ErrorType `gorm:"column:error_type"`
	ConfigurationId string             `gorm:"column:configuration_id"`
	Total           int                `gorm:"column:total"`
}
