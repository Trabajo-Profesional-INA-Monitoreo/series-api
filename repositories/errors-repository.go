package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type ErrorsRepository interface {
	AlreadyDetectedErrorForStreamWithIdAndType(uint64, string, entities.ErrorType) bool
	GetDetectedErrorForStreamWithIdAndType(streamId uint64, requestId string, errorType entities.ErrorType) entities.DetectedError
	Create(detectedError entities.DetectedError)
	Update(detectedError entities.DetectedError)
	GetErrorsPerDay(timeStart time.Time, timeEnd time.Time) []*dtos.ErrorsCountPerDayAndType
	GetErrorIndicators(timeStart time.Time, timeEnd time.Time, configId uint64) []*dtos.ErrorIndicator
}

type errorsRepository struct {
	connection *gorm.DB
}

func NewErrorsRepository(connection *gorm.DB) ErrorsRepository {
	return &errorsRepository{connection}
}

func (e errorsRepository) AlreadyDetectedErrorForStreamWithIdAndType(streamId uint64, requestId string, errorType entities.ErrorType) bool {
	var detectedError entities.DetectedError
	err := e.connection.Model(&entities.DetectedError{}).Where("stream_id = ? AND request_id = ? AND error_type = ?", streamId, requestId, errorType).First(&detectedError)
	if err.Error != nil {
		log.Errorf("Error executing AlreadyDetectedErrorForStreamWithIdAndType query: %v", err)
	}
	return detectedError.RequestId == requestId
}

func (e errorsRepository) Create(detectedError entities.DetectedError) {
	err := e.connection.Create(&detectedError)
	if err.Error != nil {
		log.Errorf("Error executing Create query: %v", err)
	}
}

func (e errorsRepository) Update(detectedError entities.DetectedError) {
	err := e.connection.Save(&detectedError)
	if err.Error != nil {
		log.Errorf("Error executing Update query: %v", err)
	}
}

func (e errorsRepository) GetDetectedErrorForStreamWithIdAndType(streamId uint64, requestId string, errorType entities.ErrorType) entities.DetectedError {
	var detectedError entities.DetectedError
	err := e.connection.Model(&entities.DetectedError{}).Where("stream_id = ? AND request_id = ? AND error_type = ?", streamId, requestId, errorType).First(&detectedError)
	if err.Error != nil {
		log.Errorf("Error executing GetDetectedErrorForStreamWithIdAndType query: %v", err)
	}
	return detectedError
}

func (e errorsRepository) GetErrorsPerDay(timeStart time.Time, timeEnd time.Time) []*dtos.ErrorsCountPerDayAndType {
	var results []*dtos.ErrorsCountPerDayAndType
	err := e.connection.Model(
		&entities.DetectedError{},
	).Select(
		"detected_errors.error_type as error_type",
		"COUNT(detected_errors.error_type) as total",
		"DATE(detected_errors.detected_date) as date",
	).Where(
		"detected_date >= ? AND detected_date <= ?", timeStart, timeEnd,
	).Group(
		"DATE(detected_errors.detected_date)",
	).Group(
		"detected_errors.error_type",
	).Find(&results)
	if err.Error != nil {
		log.Errorf("Error executing GetErrorsPerDay query: %v", err)
	}
	return results
}

func (e errorsRepository) GetErrorIndicators(timeStart time.Time, timeEnd time.Time, configId uint64) []*dtos.ErrorIndicator {
	var results []*dtos.ErrorIndicator
	err := e.connection.Model(
		&entities.DetectedError{},
	).Select(
		"detected_errors.error_type as error_type",
		"COUNT(detected_errors.error_id) as count",
	).Joins(
		"JOIN configured_streams_errors ON configured_streams_errors.detected_error_error_id=detected_errors.error_id",
	).Joins(
		"JOIN configured_streams ON configured_streams.configured_stream_id=configured_streams_errors.configured_stream_configured_stream_id",
	).Where(
		"configured_streams.configuration_id = ?", configId,
	).Where(
		"detected_date >= ? AND detected_date <= ?", timeStart, timeEnd,
	).Group(
		"detected_errors.error_type",
	).Find(&results)
	if err.Error != nil {
		log.Errorf("Error executing GetErrorIndicators query: %v", err)
	}
	return results
}
