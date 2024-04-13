package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"gorm.io/gorm"
)

type RedundancyRepository interface {
	Create(redundancy entities.Redundancy) error
}

type redundancyRepository struct {
	connection *gorm.DB
}

func (db redundancyRepository) Create(redundancy entities.Redundancy) error {
	result := db.connection.Create(&redundancy)
	return result.Error
}

func NewRedundancyRepository(connection *gorm.DB) RedundancyRepository {
	return &redundancyRepository{connection}
}
