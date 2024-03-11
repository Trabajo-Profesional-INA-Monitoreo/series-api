package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"gorm.io/gorm"
)

type ErrorsRepository interface {
	AlreadyDetectedErrorForStreamWithIdAndType(uint64, string, entities.ErrorType) bool
	GetDetectedErrorForStreamWithIdAndType(streamId uint64, requestId string, errorType entities.ErrorType) entities.DetectedError
	Create(detectedError entities.DetectedError)
	Update(detectedError entities.DetectedError)
}

type errorsRepository struct {
	connection *gorm.DB
}

func NewErrorsRepository(connection *gorm.DB) ErrorsRepository {
	return &errorsRepository{connection}
}

func (e errorsRepository) AlreadyDetectedErrorForStreamWithIdAndType(streamId uint64, requestId string, errorType entities.ErrorType) bool {
	var detectedError entities.DetectedError
	e.connection.Model(&entities.DetectedError{}).Where("stream_id = ? AND request_id = ? AND error_type = ?", streamId, requestId, errorType).First(&detectedError)
	return detectedError.RequestId == requestId
}

func (e errorsRepository) Create(detectedError entities.DetectedError) {
	e.connection.Create(&detectedError)
}

func (e errorsRepository) Update(detectedError entities.DetectedError) {
	e.connection.Save(&detectedError)
}

func (e errorsRepository) GetDetectedErrorForStreamWithIdAndType(streamId uint64, requestId string, errorType entities.ErrorType) entities.DetectedError {
	var detectedError entities.DetectedError
	e.connection.Model(&entities.DetectedError{}).Where("stream_id = ? AND request_id = ? AND error_type = ?", streamId, requestId, errorType).First(&detectedError)
	return detectedError
}
