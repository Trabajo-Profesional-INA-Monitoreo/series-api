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
	GetErrorsOfConfiguredStream(parameters *dtos.QueryParameters) (*dtos.DetectedErrorsOfStream, error)
}

type errorsService struct {
	repository repositories.ErrorMetricsRepository
}

func NewErrorsService(repository repositories.ErrorMetricsRepository) ErrorsService {
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
	return e.repository.GetRelatedStreamsToError(parameters)
}

func (e errorsService) GetErrorsOfConfiguredStream(parameters *dtos.QueryParameters) (*dtos.DetectedErrorsOfStream, error) {
	res, err := e.repository.GetErrorsOfConfiguredStream(parameters)
	if err != nil {
		return nil, err
	}
	res.ConvertToResponse()
	return res, nil
}
