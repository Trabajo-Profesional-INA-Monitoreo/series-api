package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type NotificationsRepository interface {
	GetTotalErrors(date time.Time) int64
	GetErrorsByConfigurationId(date time.Time) []*dtos.NotificationsErrorsCountPerConfigurationId
	GetErrorsByConfigurationIdAndErrorType(date time.Time) []*dtos.NotificationsErrorsCountPerType
}

type notificationsRepository struct {
	connection *gorm.DB
}

func (n *notificationsRepository) GetTotalErrors(date time.Time) int64 {
	var count int64

	err := n.connection.Raw("SELECT COUNT(detected_errors.error_id) "+
		"FROM configurations "+
		"INNER JOIN nodes ON nodes.configuration_id = configurations.configuration_id "+
		"INNER JOIN configured_streams ON (nodes.node_id = configured_streams.node_id) AND (configured_streams.configuration_id = configurations.configuration_id) "+
		"INNER JOIN configured_streams_errors ON configured_streams.configured_stream_id = configured_streams_errors.configured_stream_configured_stream_id "+
		"INNER JOIN detected_errors ON detected_errors.error_id = configured_streams_errors.detected_error_error_id "+
		"WHERE DATE(detected_errors.detected_date) = DATE(?) "+
		"AND configurations.send_notifications "+
		"AND not configurations.deleted "+
		"AND not nodes.deleted "+
		"AND not configured_streams.deleted", date.Format("2006-01-02")).Scan(&count)

	if err.Error != nil {
		log.Errorf("Error executing GetTotalErrors query: %v", err)
	}
	return count
}

func (n *notificationsRepository) GetErrorsByConfigurationId(date time.Time) []*dtos.NotificationsErrorsCountPerConfigurationId {
	var results []*dtos.NotificationsErrorsCountPerConfigurationId

	err := n.connection.Raw("SELECT configurations.configuration_id as configuration_id, configurations.name as name,COUNT(detected_errors.error_id) as total "+
		"FROM configurations "+
		"INNER JOIN nodes ON nodes.configuration_id = configurations.configuration_id "+
		"INNER JOIN configured_streams ON (nodes.node_id = configured_streams.node_id) AND (configured_streams.configuration_id = configurations.configuration_id) "+
		"INNER JOIN configured_streams_errors ON configured_streams.configured_stream_id = configured_streams_errors.configured_stream_configured_stream_id "+
		"INNER JOIN detected_errors ON detected_errors.error_id = configured_streams_errors.detected_error_error_id "+
		"WHERE DATE(detected_errors.detected_date) = DATE(?) "+
		"AND configurations.send_notifications "+
		"AND not configurations.deleted "+
		"AND not nodes.deleted "+
		"AND not configured_streams.deleted "+
		"GROUP BY configurations.configuration_id ", date.Format("2006-01-02")).Find(&results)

	if err.Error != nil {
		log.Errorf("Error executing GetErrorsByConfigurationId query: %v", err)
	}
	return results
}

func (n *notificationsRepository) GetErrorsByConfigurationIdAndErrorType(date time.Time) []*dtos.NotificationsErrorsCountPerType {
	var results []*dtos.NotificationsErrorsCountPerType

	err := n.connection.Raw("SELECT configurations.configuration_id as configuration_id, detected_errors.error_type as error_type,COUNT(detected_errors.error_id) as total "+
		"FROM configurations "+
		"INNER JOIN nodes ON nodes.configuration_id = configurations.configuration_id "+
		"INNER JOIN configured_streams ON (nodes.node_id = configured_streams.node_id) AND (configured_streams.configuration_id = configurations.configuration_id) "+
		"INNER JOIN configured_streams_errors ON configured_streams.configured_stream_id = configured_streams_errors.configured_stream_configured_stream_id "+
		"INNER JOIN detected_errors ON detected_errors.error_id = configured_streams_errors.detected_error_error_id "+
		"WHERE DATE(detected_errors.detected_date) = DATE(?) "+
		"AND configurations.send_notifications "+
		"AND not configurations.deleted "+
		"AND not nodes.deleted "+
		"AND not configured_streams.deleted "+
		"GROUP BY configurations.configuration_id, detected_errors.error_type ", date.Format("2006-01-02")).Find(&results)

	if err.Error != nil {
		log.Errorf("Error executing GetErrorsByConfigurationIdAndErrorType query: %v", err)
	}
	return results
}

func NewNotificationsRepository(connection *gorm.DB) NotificationsRepository {
	return &notificationsRepository{connection}
}
