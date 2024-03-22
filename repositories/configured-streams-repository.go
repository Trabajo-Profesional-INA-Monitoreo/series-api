package repositories

import (
	"errors"
	"fmt"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	exceptions "github.com/Trabajo-Profesional-INA-Monitoreo/series-api/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ConfiguredStreamsRepository interface {
	FindConfiguredStreamsWithCheckErrorsForStream(stream entities.Stream) []entities.ConfiguredStream
	FindConfiguredStreamById(configStreamId uint64) (entities.ConfiguredStream, error)
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
