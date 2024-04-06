package converters

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients/responses"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
)

func ConvertStreamResponseToEntity(response *responses.InaStreamResponse, streamType uint64) entities.Stream {
	return entities.Stream{
		StreamId:    uint64(response.Id),
		StationId:   uint64(response.Station.Id),
		VariableId:  uint64(response.Variable.Id),
		ProcedureId: uint64(response.Procedure.Id),
		UnitId:      uint64(response.Unit.Id),
		StreamType:  entities.StreamType(streamType),
	}
}
