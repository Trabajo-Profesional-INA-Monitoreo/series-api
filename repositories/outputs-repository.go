package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OutputsRepository interface {
	GetStreamsForOutputMetrics(configId uint64) ([]dtos.BehaviourStream, error)
}

type outputsRepository struct {
	connection *gorm.DB
}

func NewOutputsRepository(connection *gorm.DB) OutputsRepository {
	return &outputsRepository{connection}
}

func (db *outputsRepository) GetStreamsForOutputMetrics(configId uint64) ([]dtos.BehaviourStream, error) {
	var streams []dtos.BehaviourStream
	tx := db.connection.Model(
		&entities.ConfiguredStream{},
	).Select(
		"streams.stream_id as stream_id",
		"stations.alert_level as alert_level",
		"stations.evacuation_level as evacuation_level",
		"stations.low_water_level as low_water_level",
	).Joins(
		"JOIN streams ON streams.stream_id=configured_streams.stream_id",
	).Joins(
		"JOIN stations ON stations.station_id=streams.station_id",
	).Where(
		"configured_streams.configuration_id = ?", configId,
	).Where(
		"configured_streams.deleted = false",
	).Where(
		"streams.variable_id = ?", waterLevelVariable, // TODO validar
	).Where(
		"streams.procedure_id = ?", directMeasurementProcedure, // TODO validar
	).Where(
		"streams.stream_type = ?", entities.Observed, // TODO validar
	).Find(&streams)

	if tx.Error != nil {
		log.Errorf("Error executing GetStreamsForOutputMetrics query: %v", tx.Error)
		return nil, tx.Error
	}

	return streams, nil
}
