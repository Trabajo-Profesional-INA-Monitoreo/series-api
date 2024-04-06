package converters

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/clients/responses"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
)

func ConvertStationResponseToEntity(station responses.Station) entities.Station {
	return entities.Station{
		StationId:       uint64(station.Id),
		Name:            station.Name,
		Owner:           station.Owner,
		AlertLevel:      float64(station.AlertLevel),
		EvacuationLevel: station.EvacuationLevel,
		LowWaterLevel:   station.LowWaterLevel,
	}
}
