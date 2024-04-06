package converters

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients/responses"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
)

func ConvertUnitResponseToEntity(unit responses.Unit) entities.Unit {
	return entities.Unit{UnitId: uint64(unit.Id), Name: unit.Name}
}
