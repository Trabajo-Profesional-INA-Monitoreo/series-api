package services

import (
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/dtos"
	"github.com/Trabajo-Profesional-INA-Monitoreo/series-api/repositories"
	"time"
)

type ErrorsService interface {
	GetErrorsPerDay(timeStart time.Time, timeEnd time.Time) []*dtos.ErrorsCountPerDayAndType
}

type errorsService struct {
	repository repositories.ErrorsRepository
}

func NewErrorsService(repository repositories.ErrorsRepository) ErrorsService {
	return &errorsService{repository: repository}
}

func (e errorsService) GetErrorsPerDay(timeStart time.Time, timeEnd time.Time) []*dtos.ErrorsCountPerDayAndType {
	errors := e.repository.GetErrorsPerDay(timeStart, timeEnd)
	for _, count := range errors {
		count.ConvertToResponse()
	}
	return errors
}
