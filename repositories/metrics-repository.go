package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MetricsRepository interface {
	Create(metric entities.ConfiguredMetric) error
	GetByConfiguredStreamId(id uint64) *[]entities.Metric
}

type metricsRepository struct {
	connection *gorm.DB
}

func (db metricsRepository) GetByConfiguredStreamId(configuredStreamId uint64) *[]entities.Metric {
	var metrics *[]entities.Metric

	tx := db.connection.Model(
		&entities.ConfiguredMetric{},
	).Select(
		"configured_metrics.metric_id as metric_id",
	).Where(
		"configured_metrics.configured_stream_id = ?", configuredStreamId,
	).Scan(&metrics)

	if tx.Error != nil {
		log.Errorf("Error executing GetMetrics query: %v", tx.Error)
	}

	log.Debugf("Get metrics query result: %v", metrics)
	return metrics
}

func (db metricsRepository) Create(metric entities.ConfiguredMetric) error {
	result := db.connection.Create(&metric)
	return result.Error
}

func NewMetricsRepository(connection *gorm.DB) MetricsRepository {
	return &metricsRepository{connection}
}
