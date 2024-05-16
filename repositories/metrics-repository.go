package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MetricsRepository interface {
	Create(metric entities.ConfiguredMetric) error
	GetByConfiguredStreamId(id uint64) *[]entities.Metric
	DeleteMetricsNotIncludedInNewConfig(id uint64, metrics []entities.Metric)
}

type metricsRepository struct {
	connection *gorm.DB
}

func (db metricsRepository) DeleteMetricsNotIncludedInNewConfig(id uint64, metrics []entities.Metric) {
	tx := db.connection.Where(
		"configured_stream_id = ?", id,
	)

	if len(metrics) != 0 {
		tx = tx.Where("metric_id NOT IN ?", metrics)
	}
	tx.Delete(&entities.ConfiguredMetric{})

	if tx.Error != nil {
		log.Errorf("Error executing DeleteMetricsNotIncludedInNewConfig query: %v", tx.Error)
	}
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
