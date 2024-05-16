package converters

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients/responses"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
)

func ConvertVariableResponseToEntity(variable responses.Variable) entities.Variable {
	return entities.Variable{
		VariableId: uint64(variable.Id),
		Name:       variable.Name,
	}
}
