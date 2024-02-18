package repositories

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/entities"
	"gorm.io/gorm"
)

type ErrorsRepository interface {
	AlreadyDetectedErrorWithIdAndType(string, entities.ErrorType) bool
	Save(detectedError entities.DetectedError)
}

type errorsRepository struct {
	connection *gorm.DB
}

func NewErrorsRepository(connection *gorm.DB) ErrorsRepository {
	return &errorsRepository{connection}
}

func (e errorsRepository) AlreadyDetectedErrorWithIdAndType(requestId string, errorType entities.ErrorType) bool {
	var detectedError entities.DetectedError
	e.connection.Model(&entities.DetectedError{}).Where("request_id = ? AND error_type = ?", requestId, errorType).First(&detectedError)
	return detectedError.RequestId == requestId
}

func (e errorsRepository) Save(detectedError entities.DetectedError) {
	e.connection.Create(&detectedError)
}
