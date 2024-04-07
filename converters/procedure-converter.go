package converters

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients/responses"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
)

func ConvertProcedureResponseToEntity(procedure responses.Procedure) entities.Procedure {
	return entities.Procedure{
		ProcedureId: uint64(procedure.Id),
		Name:        procedure.Description,
	}
}
