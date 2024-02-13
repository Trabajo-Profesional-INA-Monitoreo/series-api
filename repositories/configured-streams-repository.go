package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"gorm.io/gorm"
)

type ConfiguredStreamsRepository interface {
	FindConfiguredStreamsForStream(stream entities.Stream) []entities.ConfiguredStream
}

type configuredStreamsRepository struct {
	connection *gorm.DB
}

func NewConfiguredStreamsRepository(connection *gorm.DB) ConfiguredStreamsRepository {
	return &configuredStreamsRepository{connection}
}

func (c configuredStreamsRepository) FindConfiguredStreamsForStream(stream entities.Stream) []entities.ConfiguredStream {
	//TODO implement me
	panic("implement me")
}
