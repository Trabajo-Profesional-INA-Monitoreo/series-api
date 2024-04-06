package repositories

import (
	"errors"
	"fmt"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	exceptions "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ConfiguredStreamsRepository interface {
	FindConfiguredStreamsWithCheckErrorsForStream(stream entities.Stream) []entities.ConfiguredStream
	FindConfiguredStreamById(configStreamId uint64) (entities.ConfiguredStream, error)
	Create(e *entities.ConfiguredStream) error
	FindConfiguredStreamsByNodeId(nodeId uint64, configurationId string) *[]dtos.ConfiguredStream
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

func (db configuredStreamsRepository) Create(configuredStream *entities.ConfiguredStream) error {
	result := db.connection.Create(&configuredStream)
	return result.Error
}

func (db configuredStreamsRepository) FindConfiguredStreamsByNodeId(nodeId uint64, configurationId string) *[]dtos.ConfiguredStream {
	var configuredStream []dtos.ConfiguredStream

	result := db.connection.Model(
		&entities.ConfiguredStream{},
	).Select(
		"configured_streams.stream_id ",
		"streams.stream_type ",
		"configured_streams.update_frequency",
		"configured_streams.check_errors",
		"configured_streams.normal_upper_threshold",
		"configured_streams.normal_lower_threshold",
		"configured_streams.calibration_id",
	).Where("node_id = ? AND configuration_id = ?", nodeId, configurationId).Joins("JOIN streams ON streams.stream_id = configured_streams.stream_id ").Scan(&configuredStream)

	if result.RowsAffected == 0 {
		return nil
	}

	log.Debugf("Get configurations query result: %v", configuredStream)
	return &configuredStream
}
