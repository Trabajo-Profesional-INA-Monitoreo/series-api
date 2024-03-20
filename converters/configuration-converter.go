package converters

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
)

func ConvertDtoToConfiguration(dto dtos.Configuration) *entities.Configuration {
	return &entities.Configuration{Name: dto.Name, Id: dto.Id}
}
