package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RedundancyRepository interface {
	Create(redundancy entities.Redundancy) error
	GetByConfiguredStreamId(id uint64) *[]uint64
}

type redundancyRepository struct {
	connection *gorm.DB
}

func (db redundancyRepository) GetByConfiguredStreamId(configuredStreamId uint64) *[]uint64 {
	var redundancies *[]uint64

	tx := db.connection.Model(
		&entities.Redundancy{},
	).Select(
		"redundancies.redundancy_id as redundancy_id",
	).Where(
		"redundancies.configured_stream_id = ?", configuredStreamId,
	).Scan(&redundancies)

	if tx.Error != nil {
		log.Errorf("Error executing GetRedundancies query: %v", tx.Error)
	}

	log.Debugf("Get redundancies query result: %v", redundancies)
	return redundancies
}

func (db redundancyRepository) Create(redundancy entities.Redundancy) error {
	result := db.connection.Create(&redundancy)
	return result.Error
}

func NewRedundancyRepository(connection *gorm.DB) RedundancyRepository {
	return &redundancyRepository{connection}
}
