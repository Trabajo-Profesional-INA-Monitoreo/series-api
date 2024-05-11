package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type FilterRepository interface {
	GetProcedures(configId uint64) []dtos.FilterValue
	GetStations(configId uint64) []dtos.FilterValue
	GetVariables(configId uint64) []dtos.FilterValue
	GetNodes(configId uint64) []dtos.FilterValue
}

type filterRepository struct {
	connection *gorm.DB
}

func (f filterRepository) GetNodes(configId uint64) []dtos.FilterValue {
	var filters []dtos.FilterValue

	f.connection.Model(
		&entities.Node{},
	).Select(
		"nodes.name as name, nodes.node_id as id",
	).Where("nodes.configuration_id = ?", configId).Scan(&filters)

	log.Debugf("Get nodes query result: %v", filters)
	return filters
}

func (f filterRepository) GetProcedures(configId uint64) []dtos.FilterValue {
	var filters []dtos.FilterValue

	f.connection.Model(
		&entities.ConfiguredStream{},
	).Select(
		"distinct(procedures.name) as name, procedures.procedure_id as id",
	).Joins(
		"JOIN streams ON configured_streams.stream_id = streams.stream_id",
	).Joins(
		"JOIN procedures ON streams.procedure_id = procedures.procedure_id",
	).Where(
		"configured_streams.configuration_id = ?", configId,
	).Scan(&filters)

	log.Debugf("Get procedures query result: %v", filters)
	return filters
}

func (f filterRepository) GetStations(configId uint64) []dtos.FilterValue {
	var filters []dtos.FilterValue

	f.connection.Model(
		&entities.ConfiguredStream{},
	).Select(
		"distinct(stations.name) as name, stations.station_id as id",
	).Joins(
		"JOIN streams ON configured_streams.stream_id = streams.stream_id",
	).Joins(
		"JOIN stations ON streams.station_id = stations.station_id",
	).Where(
		"configured_streams.configuration_id = ?", configId,
	).Scan(&filters)

	log.Debugf("Get stations query result: %v", filters)
	return filters
}

func (f filterRepository) GetVariables(configId uint64) []dtos.FilterValue {
	var filters []dtos.FilterValue

	f.connection.Model(
		&entities.ConfiguredStream{},
	).Select(
		"distinct(variables.name) as name, variables.variable_id as id",
	).Joins(
		"JOIN streams ON configured_streams.stream_id = streams.stream_id",
	).Joins(
		"JOIN variables ON streams.variable_id = variables.variable_id",
	).Where(
		"configured_streams.configuration_id = ?", configId,
	).Scan(&filters)

	log.Debugf("Get variable query result: %v", filters)
	return filters
}

func NewFilterRepository(connection *gorm.DB) FilterRepository {
	return &filterRepository{connection}
}
