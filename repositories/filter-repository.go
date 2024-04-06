package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type FilterRepository interface {
	GetProcedures() []dtos.ProcedureFilter
	GetStations() []dtos.StationFilter
	GetVariables() []dtos.VariableFilter
}

type filterRepository struct {
	connection *gorm.DB
}

func (f filterRepository) GetProcedures() []dtos.ProcedureFilter {
	var procedures []dtos.ProcedureFilter

	f.connection.Model(
		&entities.Procedure{},
	).Select(
		"procedures.name as name, procedures.procedure_id as id",
	).Scan(&procedures)

	log.Debugf("Get procedures query result: %v", procedures)
	return procedures
}

func (f filterRepository) GetStations() []dtos.StationFilter {
	var stations []dtos.StationFilter

	f.connection.Model(
		&entities.Station{},
	).Select(
		"stations.name as name, stations.station_id as id",
	).Scan(&stations)

	log.Debugf("Get stations query result: %v", stations)
	return stations
}

func (f filterRepository) GetVariables() []dtos.VariableFilter {
	var variable []dtos.VariableFilter

	f.connection.Model(
		&entities.Variable{},
	).Select(
		"variables.name as name, variables.variable_id as id",
	).Scan(&variable)

	log.Debugf("Get variable query result: %v", variable)
	return variable
}

func NewFilterRepository(connection *gorm.DB) FilterRepository {
	return &filterRepository{connection}
}
