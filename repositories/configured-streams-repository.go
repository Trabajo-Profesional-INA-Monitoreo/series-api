package repositories

import (
	"errors"
	"fmt"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	exceptions "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type ConfiguredStreamsRepository interface {
	FindConfiguredStreamsWithCheckErrorsForStream(stream entities.Stream) []entities.ConfiguredStream
	FindConfiguredStreamById(configStreamId uint64) (entities.ConfiguredStream, error)
	CountErrorOfConfigurations(ids []uint64, parameters *dtos.QueryParameters) ([]dtos.ErrorsPerConfigStream, error)
}

type configuredStreamsRepository struct {
	connection *gorm.DB
}

func NewConfiguredStreamsRepository(connection *gorm.DB) ConfiguredStreamsRepository {
	return &configuredStreamsRepository{connection}
}

func (db configuredStreamsRepository) FindConfiguredStreamsWithCheckErrorsForStream(stream entities.Stream) []entities.ConfiguredStream {
	var configured []entities.ConfiguredStream

	db.connection.Model(
		&entities.ConfiguredStream{},
	).Where("stream_id = ? AND check_errors = true", stream.StreamId).Find(&configured)

	return configured
}

func (db configuredStreamsRepository) FindConfiguredStreamById(configStreamId uint64) (entities.ConfiguredStream, error) {
	var configured entities.ConfiguredStream

	result := db.connection.Model(
		&entities.ConfiguredStream{},
	).Preload("Metrics").Where("configured_stream_id = ?", configStreamId).Find(&configured)

	if result.Error != nil {
		log.Errorf("Error executing FindConfiguredStreamById query: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return configured, errors.Join(exceptions.NewNotFound(), fmt.Errorf("configured stream with id %v not found", configStreamId))
	}

	return configured, nil
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
