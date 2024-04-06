package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	"time"
)

type ErrorsService interface {
	GetErrorsPerDay(timeStart time.Time, timeEnd time.Time, configId uint64) []*dtos.ErrorsCountPerDayAndType
	GetErrorIndicators(timeStart time.Time, timeEnd time.Time, configId uint64) []*dtos.ErrorIndicator
	GetRelatedStreams(parameters *dtos.QueryParameters) ([]dtos.ErrorRelatedStream, error)
}

type errorsService struct {
	repository repositories.ErrorsRepository
}

func NewErrorsService(repository repositories.ErrorsRepository) ErrorsService {
	return &errorsService{repository: repository}
}

func (e errorsService) GetErrorsPerDay(timeStart time.Time, timeEnd time.Time, configId uint64) []*dtos.ErrorsCountPerDayAndType {
	errors := e.repository.GetErrorsPerDay(timeStart, timeEnd, configId)
	for _, count := range errors {
		count.ConvertToResponse()
	}
	return errors
}

func (e errorsService) GetErrorIndicators(timeStart time.Time, timeEnd time.Time, configId uint64) []*dtos.ErrorIndicator {
	errors := e.repository.GetErrorIndicators(timeStart, timeEnd, configId)
	for _, count := range errors {
		count.ConvertToResponse()
	}
	return errors
}

func (e errorsService) GetRelatedStreams(parameters *dtos.QueryParameters) ([]dtos.ErrorRelatedStream, error) {
	related, err := e.repository.GetRelatedStreamsToError(parameters)
	if err != nil {
		return nil, err
	}
	return related, nil
}
