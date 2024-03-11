package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"gorm.io/gorm"
)

type ConfiguredStreamsRepository interface {
	FindConfiguredStreamsWithCheckErrorsForStream(stream entities.Stream) []entities.ConfiguredStream
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
