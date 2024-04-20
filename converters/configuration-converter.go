package converters

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
)

func ConvertDtoToCreateConfiguration(dto dtos.CreateConfiguration) *entities.Configuration {
	return &entities.Configuration{Name: dto.Name, Deleted: false, SendNotifications: dto.SendNotifications}
}

func ConvertDtoToConfiguration(dto dtos.Configuration) *entities.Configuration {
	return &entities.Configuration{Name: dto.Name, ConfigurationId: dto.Id, SendNotifications: dto.SendNotifications}
}
