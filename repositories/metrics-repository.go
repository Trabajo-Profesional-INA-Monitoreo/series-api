package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"gorm.io/gorm"
)

type MetricsRepository interface {
	Create(metric entities.ConfiguredMetric) error
}

type metricsRepository struct {
	connection *gorm.DB
}

func (db metricsRepository) Create(metric entities.ConfiguredMetric) error {
	result := db.connection.Create(&metric)
	return result.Error
}

func NewMetricsRepository(connection *gorm.DB) MetricsRepository {
	return &metricsRepository{connection}
}
