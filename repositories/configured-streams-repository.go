package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	exceptions "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MetricsConfiguredStreamsRepository interface {
	FindConfiguredStreamById(configStreamId uint64) (entities.ConfiguredStream, error)
	CountErrorOfConfigurations(ids []uint64, parameters *dtos.QueryParameters) ([]dtos.ErrorsPerConfigStream, error)
}

type DetectionConfiguredStreamsRepository interface {
	FindConfiguredStreamsWithCheckErrorsForStream(stream entities.Stream) []entities.ConfiguredStream
}

type ManagerConfiguredStreamsRepository interface {
	Create(e *entities.ConfiguredStream) (uint64, error)
	FindConfiguredStreamsByNodeId(nodeId uint64, configurationId uint64) *[]*dtos.ConfiguredStream
	Update(e *entities.ConfiguredStream) error
	MarkAsDeletedOldConfiguredStreams(configId uint64, newConfigStreamIds []uint64)
	DeleteByConfig(configId uint64)
}

type ConfiguredStreamsRepository interface {
	MetricsConfiguredStreamsRepository
	DetectionConfiguredStreamsRepository
	ManagerConfiguredStreamsRepository
}

type configuredStreamsRepository struct {
	connection *gorm.DB
}

func (db configuredStreamsRepository) DeleteByConfig(configId uint64) {
	db.connection.Where(
		"configured_metrics.configured_stream_id IN "+
			"(SELECT configured_stream_id FROM configured_streams "+
			"WHERE configured_streams.configuration_id = ? "+
			"AND configured_streams.deleted = true)", configId,
	).Delete(&entities.ConfiguredMetric{})

	db.connection.Where(
		"redundancies.configured_stream_id IN "+
			"(SELECT configured_stream_id FROM configured_streams "+
			"WHERE configured_streams.configuration_id = ? "+
			"AND configured_streams.deleted = true)", configId,
	).Delete(&entities.Redundancy{})

	db.connection.Exec("DELETE FROM configured_streams_errors WHERE "+
		"configured_streams_errors.configured_stream_configured_stream_id IN "+
		"(SELECT configured_stream_id FROM configured_streams "+
		"WHERE configured_streams.configuration_id = ? "+
		"AND configured_streams.deleted = true)", configId,
	)

	db.connection.Where(
		"configured_streams.deleted = true",
	).Where("configured_streams.configuration_id = ?", configId).Delete(&entities.ConfiguredStream{})
}

func (db configuredStreamsRepository) MarkAsDeletedOldConfiguredStreams(configId uint64, newConfigStreamIds []uint64) {
	tx := db.connection.Model(
		&entities.ConfiguredStream{},
	).Where(
		"configured_streams.configuration_id = ?", configId,
	)
	if newConfigStreamIds != nil && len(newConfigStreamIds) != 0 {
		tx = tx.Where(
			"configured_streams.configured_stream_id NOT IN ?", newConfigStreamIds,
		)
	}
	tx.Update("deleted", true)
}

func NewConfiguredStreamsRepository(connection *gorm.DB) ConfiguredStreamsRepository {
	return &configuredStreamsRepository{connection}
}

func (db configuredStreamsRepository) FindConfiguredStreamsWithCheckErrorsForStream(stream entities.Stream) []entities.ConfiguredStream {
	var configured []entities.ConfiguredStream

	db.connection.Model(
		&entities.ConfiguredStream{},
	).Preload("Configuration").Where("stream_id = ? AND check_errors = true AND deleted = false", stream.StreamId).Find(&configured)

	return configured
}

func (db configuredStreamsRepository) FindConfiguredStreamById(configStreamId uint64) (entities.ConfiguredStream, error) {
	var configured entities.ConfiguredStream

	result := db.connection.Model(
		&entities.ConfiguredStream{},
	).Preload("Metrics").Where("configured_stream_id = ? AND deleted = false", configStreamId).Find(&configured)

	if result.Error != nil {
		log.Errorf("Error executing FindConfiguredStreamById query: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return configured, errors.Join(exceptions.NewNotFound(), fmt.Errorf("configured stream with id %v not found", configStreamId))
	}

	return configured, nil
}

func (db configuredStreamsRepository) Create(configuredStream *entities.ConfiguredStream) (uint64, error) {
	result := db.connection.Create(&configuredStream)
	return configuredStream.ConfiguredStreamId, result.Error
}

func (db configuredStreamsRepository) FindConfiguredStreamsByNodeId(nodeId uint64, configurationId uint64) *[]*dtos.ConfiguredStream {
	var configuredStream *[]*dtos.ConfiguredStream

	result := db.connection.Model(
		&entities.ConfiguredStream{},
	).Select(
		"configured_streams.configured_stream_id",
		"configured_streams.stream_id ",
		"streams.stream_type ",
		"configured_streams.update_frequency",
		"configured_streams.check_errors",
		"configured_streams.normal_upper_threshold",
		"configured_streams.normal_lower_threshold",
		"configured_streams.calibration_id",
		"configured_streams.forecasted_range_hours",
		"configured_streams.observed_related_stream_id",
	).Where(
		"node_id = ? AND configuration_id = ?", nodeId, configurationId,
	).Where("configured_streams.deleted = false").Joins("JOIN streams ON streams.stream_id = configured_streams.stream_id ").Scan(&configuredStream)

	if result.RowsAffected == 0 {
		return nil
	}

	log.Debugf("Get configurations query result: %v", configuredStream)
	return configuredStream
}

func (db configuredStreamsRepository) Update(configuredStream *entities.ConfiguredStream) error {
	result := db.connection.Save(&configuredStream)
	return result.Error
}

func (db configuredStreamsRepository) CountErrorOfConfigurations(ids []uint64, parameters *dtos.QueryParameters) ([]dtos.ErrorsPerConfigStream, error) {
	timeStart := parameters.Get("timeStart").(time.Time)
	timeEnd := parameters.Get("timeEnd").(time.Time)
	var results []dtos.ErrorsPerConfigStream
	tx := db.connection.Model(
		&entities.ConfiguredStream{},
	).Select(
		"configured_streams.configured_stream_id as configured_stream_id",
		"COUNT(configured_streams_errors.detected_error_error_id) as errors_count",
	).Joins(
		"JOIN configured_streams_errors ON configured_streams_errors.configured_stream_configured_stream_id=configured_streams.configured_stream_id",
	).Joins(
		"JOIN detected_errors ON detected_errors.error_id=configured_streams_errors.detected_error_error_id",
	).Where(
		"configured_streams.configured_stream_id IN ?", ids,
	).Where(
		"detected_errors.detected_date >= ? AND detected_errors.detected_date <= ?", timeStart, timeEnd,
	).Group("configured_streams.configured_stream_id").Find(&results)

	if tx.Error != nil {
		log.Errorf("Error executing CountErrorOfConfigurations query: %v", tx.Error)
		return nil, tx.Error
	}
	return results, nil
}
